package game

import (
	"encoding/json"
	"fmt"
)

const GAME_WIDTH = 100

type PlayerId int

const (
	NoPlayerId PlayerId = iota - 1
	OnThePlay
	OnTheDraw
)

func (id PlayerId) OpponentId() PlayerId {
	switch id {
	case OnThePlay:
		return OnTheDraw
	case OnTheDraw:
		return OnThePlay
	default:
		panic("bad id")
	}
	panic("bad logic")
}

type Game struct {
	// Players are sometimes referred to by index, 0 or 1.
	// Player 0 is the player who plays first.
	Players [2]*Player
	Phase   Phase

	// Which turn of the game is. This starts at 0.
	// In general, it is player (turn % 2)'s turn.
	Turn int

	// The id of the player with priority
	PriorityId PlayerId

	// The PermanentId that will be assigned to the next permanent that enters play
	NextPermanentId PermanentId

	// Permanents contains all permanents in play.
	Permanents map[PermanentId]*Permanent

	/*
		The effect the Priority Player must decide on
		Currently, the only ChoiceEffect is how to pay for Daze
	*/
	ChoiceEffect *Effect

	/*
		Some actions go on the stack and can be responded to before they resolve
		https://mtg.gamepedia.com/Stack#Actions
	*/

	// StackObjects contains all objects on the stack
	StackObjects map[StackObjectId]*StackObject
	// The stack holds IDs only to prevent circular refs in serialization.
	Stack []StackObjectId
	// The StackObjectId that will be assigned to the next object that gets put on the stack.
	NextStackObjectId StackObjectId

	// True if the acting player passed priority after putting a spell or ability on the stack.
	ActorPassedOnStack bool
}

//go:generate stringer -type=Phase
type Phase int

// Phase encompasses both "phases" and "steps" as per:
// https://mtg.gamepedia.com/Turn_structure
// For example, declaring attackers is a step within the combat phase in the official
// rules. Here it is just treated as another Phase.

const (
	UntapStep Phase = iota
	Upkeep
	Draw
	Main1
	DeclareAttackers
	DeclareBlockers
	CombatDamage
	Main2
)

func NewGame(deckToPlay *Deck, deckToDraw *Deck) *Game {
	players := [2]*Player{
		NewPlayer(deckToPlay, OnThePlay),
		NewPlayer(deckToDraw, OnTheDraw),
	}
	g := &Game{
		Players:           players,
		Phase:             Main1,
		Turn:              0,
		PriorityId:        OnThePlay,
		NextPermanentId:   PermanentId(1),
		NextStackObjectId: StackObjectId(1),
		Permanents:        make(map[PermanentId]*Permanent),
		Stack:             []StackObjectId{},
		StackObjects:      make(map[StackObjectId]*StackObject),
	}

	players[0].game = g
	players[1].game = g

	return g
}

// Player(id) returns the player with the provided id.
func (g *Game) Player(id PlayerId) *Player {
	return g.Players[id]
}

func (g *Game) Priority() *Player {
	return g.Players[g.PriorityId]
}

func (g *Game) Actions(forHuman bool) []*Action {
	actions := []*Action{}
	// forHuman = false

	// TODO maybe some other data stucture beside ChoiceEffect - a pointer to the action on stack instead?
	// Currently handles Daze, Scry effects, and Ponder
	if g.ChoiceEffect != nil {
		return g.Priority().OptionsForChoiceEffect(g.ChoiceEffect)
	}

	if len(g.Stack) > 0 {
		actions = append(actions, &Action{
			Type: PassPriority,
		})
		stackObjectId := g.Stack[len(g.Stack)-1]
		stackObject := g.StackObject(stackObjectId)

		if g.PriorityId == stackObject.Player && g.ActorPassedOnStack {
			return actions
		}
		actions = append(actions, g.Priority().PlayActions(false, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(false, forHuman)...)
		return actions
	}

	currentPlayerIsActing := g.PriorityId == g.AttackerId()
	switch g.Phase {
	case UntapStep:
		return append(actions, g.Priority().PassAction())
	case Upkeep:
		fallthrough
	case Draw:
		actions = append(actions, g.Priority().PlayActions(false, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(false, forHuman)...)
		return addManaAndPassActions(forHuman, g, actions)
	case Main1:
		fallthrough
	case Main2:
		actions = append(actions, g.Priority().PlayActions(currentPlayerIsActing, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(currentPlayerIsActing, forHuman)...)
		return addManaAndPassActions(forHuman, g, actions)
	case DeclareAttackers:
		return append(g.Priority().AttackActions(), g.Priority().PassAction())
	case DeclareBlockers:
		return append(g.Priority().BlockActions(), g.Priority().PassAction())
	case CombatDamage:
		actions = append(actions, g.Priority().PlayActions(false, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(false, forHuman)...)
		return addManaAndPassActions(forHuman, g, actions)
	default:
		panic("unhandled phase")
	}
}

func addManaAndPassActions(forHuman bool, g *Game, actions []*Action) []*Action {
	if forHuman && len(actions) <= 1 {
		actions = appendPassAction(g, actions)
		return actions
	} else {
		actions = append(actions, g.Priority().ManaActions()...)
		actions = appendPassAction(g, actions)
		return actions
	}
}

func appendPassAction(g *Game, actions []*Action) []*Action {
	if g.PriorityId == g.AttackerId() {
		return append(actions, &Action{Type: PassPriority})
	} else {
		return append(actions, g.Priority().PassAction())
	}
}

func (g *Game) AttackerId() PlayerId {
	return PlayerId((g.Turn) % 2)
}

func (g *Game) Attacker() *Player {
	return g.Players[g.AttackerId()]
}

func (g *Game) DefenderId() PlayerId {
	return g.AttackerId().OpponentId()
}

func (g *Game) Defender() *Player {
	return g.Players[g.DefenderId()]
}

func (g *Game) HandleCombatDamage() {
	for _, attacker := range g.Attacker().GetBoard() {
		if attacker.Attacking {
			damage := attacker.Power()
			if damage < 0 {
				damage = 0
			}

			if len(attacker.DamageOrder) > 0 {
				// Deal damage to blockers
				for _, blocker := range attacker.GetDamageOrder() {
					attacker.Damage += blocker.Power()
					if damage == 0 {
						continue
					}
					remaining := blocker.Toughness() - blocker.Damage
					if remaining > damage {
						blocker.Damage += damage
						damage = 0
					} else {
						g.Defender().SendToGraveyard(blocker)
						damage -= remaining
					}
				}
			}

			if len(attacker.DamageOrder) == 0 || attacker.Trample() {
				// Deal damage to the defending player
				g.Defender().DealDamage(damage)
				attacker.DidDealDamage(damage)
			}

			if attacker.Damage >= attacker.Toughness() {
				g.Attacker().SendToGraveyard(attacker)
			}
		}
	}
}

// Creatures() returns the creatures in play.
func (g *Game) Creatures() []*Permanent {
	answer := []*Permanent{}
	for _, player := range g.Players {
		answer = append(answer, player.Creatures()...)
	}
	return answer
}

// Lands() returns the lands in play.
func (g *Game) Lands() []*Permanent {
	answer := []*Permanent{}
	for _, player := range g.Players {
		answer = append(answer, player.Lands()...)
	}
	return answer
}

func (g *Game) nextPhase() {
	for _, p := range g.Players {
		p.EndPhase()
	}

	g.ActorPassedOnStack = false
	switch g.Phase {
	case UntapStep:
		g.Attacker().Untap()
		g.Phase = Upkeep
	case Upkeep:
		for _, c := range g.Attacker().GetBoard() {
			if c.BeginningOfYourUpkeepEffect != nil {
				g.Attacker().ResolveEffect(c.BeginningOfYourUpkeepEffect, c)
			}
		}
		g.Phase = Draw
	case Draw:
		g.Priority().Draw()
		g.Phase = Main1
	case Main1:
		g.Phase = DeclareAttackers
	case DeclareAttackers:
		g.Phase = DeclareBlockers
		g.PriorityId = g.DefenderId()
	case DeclareBlockers:
		g.Phase = CombatDamage
		g.PriorityId = g.AttackerId()
	case CombatDamage:
		g.HandleCombatDamage()
		g.Attacker().EndCombat()
		g.Defender().EndCombat()
		g.Phase = Main2
	case Main2:
		// End the turn
		for _, p := range g.Players {
			p.EndTurn()
		}
		g.Phase = UntapStep
		g.Turn++
		g.PriorityId = g.PriorityId.OpponentId()
	}
}

func (g *Game) TakeAction(action *Action) {
	if g.IsOver() {
		panic("cannot take action when the game is over")
	}
	if action.Type == MakeChoice {
		if action.ShouldSwitchPriority {
			g.PriorityId = g.PriorityId.OpponentId()
		}
		g.Priority().ResolveEffect(action.AfterEffect, nil)
		g.ChoiceEffect = nil
		return
	}

	if action.Type == PassPriority {
		g.ActorPassedOnStack = true
		g.PriorityId = g.PriorityId.OpponentId()
		if len(g.Stack) > 0 {
			stackObject := g.StackObject(g.Stack[len(g.Stack)-1])
			if g.PriorityId == stackObject.Player {
				g.ActorPassedOnStack = false
				g.Stack = g.Stack[:len(g.Stack)-1]
				if stackObject.Type == Play {
					g.Player(stackObject.Player).ResolveSpell(stackObject)
				} else if stackObject.Type == Activate {
					g.Player(stackObject.Player).ResolveActivatedAbility(stackObject)
				} else if stackObject.Type == EntersTheBattlefieldEffect {
					for _, perm := range g.Priority().GetBoard() {
						if perm.Card == stackObject.Card {
							effect := UpdatedEffectForStackObject(stackObject, stackObject.Card.EntersTheBattlefieldEffect)
							g.Priority().ResolveEffect(effect, perm)
							break
						}
					}
				}
				delete(g.StackObjects, stackObject.Id)
			}
		}
		return
	}
	if action.Type == Pass {
		g.PriorityId = g.AttackerId()
		g.nextPhase()
		return
	}

	if action.Type == UseForMana {
		g.Permanent(action.Source).UseForMana()
		return
	}

	switch g.Phase {

	case Upkeep:
		fallthrough
	case Draw:
		fallthrough
	case Main1:
		fallthrough
	case Main2:
		if action.Type == Play {
			if action.Card.IsLand() {
				g.Priority().PlayLand(action)
			} else {
				g.Priority().PayCostsAndPutSpellOnStack(action)
			}
		} else if action.Type == Activate {
			g.Priority().PayCostsAndPutAbilityOnStack(action)
		} else {
			panic("expected a play, activate, declare attack, or pass during main phase")
		}

	case DeclareAttackers:
		if action.Type != Attack {
			panic("expected an attack or a pass during DeclareAttackers")
		}
		creature := g.Permanent(action.With)
		creature.Attacking = true
		creature.Tapped = true

	case DeclareBlockers:
		if action.Type != Block {
			panic("expected a block or a pass during DeclareBlockers")
		}
		creature := g.Permanent(action.With)
		creature.Blocking = action.Target
		perm := g.Permanent(action.Target)
		perm.DamageOrder = append(perm.DamageOrder, creature.Id)

	case CombatDamage:
		if action.Type == Play {
			g.Priority().PayCostsAndPutSpellOnStack(action)
		} else if action.Type == Activate {
			g.Priority().PayCostsAndPutAbilityOnStack(action)
		} else {
			panic("expected a play or activate during CombatDamage")
		}

	default:
		panic("unhandled phase")
	}
}

// Removes targetSpell from the stack, as in when Counterspelled.
func (g *Game) RemoveSpellFromStack(targetSpell StackObjectId) {
	newStack := []StackObjectId{}
	for _, spellAction := range g.Stack {
		if spellAction != targetSpell {
			newStack = append(newStack, spellAction)
		}
	}
	if len(newStack) == len(g.Stack) {
		// fmt.Println("This should be fine, it means a Counterspell's target was countered.")
	}
	g.Stack = newStack
	delete(g.StackObjects, targetSpell)
}

func (g *Game) IsOver() bool {
	return g.Attacker().Lost() || g.Defender().Lost()
}

// All permanents added to the game should be created via newPermanent.
// This assigns a unique id to the permanent and activates any coming-into-play
// effects.
func (g *Game) newPermanent(card *Card, ownerId PlayerId, stackObjectId StackObjectId, addToBoard bool) *Permanent {
	perm := &Permanent{
		Card:       card,
		Owner:      ownerId,
		TurnPlayed: g.Turn,
		Id:         g.NextPermanentId,
		game:       g,
	}
	owner := g.Player(ownerId)
	if addToBoard {
		g.Permanents[g.NextPermanentId] = perm
		owner.Board = append(owner.Board, perm.Id)
		g.NextPermanentId++
		perm.HandleEnterTheBattlefield(stackObjectId)
	}
	return perm
}

// removePermanent does nothing if the permanent has already been removed.
func (g *Game) removePermanent(id PermanentId) {
	delete(g.Permanents, id)
}

func (g *Game) Permanent(id PermanentId) *Permanent {
	if id == PermanentId(0) {
		panic("0 is not a valid PermanentId")
	}
	return g.Permanents[id]
}

func (g *Game) GetPermanents(ids []PermanentId) []*Permanent {
	answer := []*Permanent{}
	for _, id := range ids {
		answer = append(answer, g.Permanent(id))
	}
	return answer
}

// All objects must be put on the Stack via this method to get a unique ID.
func (g *Game) AddToStack(stackObject *StackObject) {
	stackObject.Id = g.NextStackObjectId
	g.StackObjects[stackObject.Id] = stackObject
	g.NextStackObjectId++
	g.Stack = append(g.Stack, stackObject.Id)
}

func (g *Game) StackObject(id StackObjectId) *StackObject {
	if id == NoStackObjectId {
		panic("NoStackObjectId is not a valid StackObjectId")
	}
	obj, ok := g.StackObjects[id]
	if !ok {
		// it's OK, spell must have been otherwise countered
	}
	return obj
}

func (g *Game) GetStack() []*StackObject {
	answer := []*StackObject{}
	for _, id := range g.Stack {
		answer = append(answer, g.StackObject(id))
	}
	return answer
}

func (g *Game) Print() {
	gameWidth := GAME_WIDTH
	printBorder(gameWidth)
	g.Players[1].Print(1, false, gameWidth)
	printMiddleLine(gameWidth)
	g.Players[0].Print(0, false, gameWidth)
	printBorder(gameWidth)
}

func printBorder(gameWidth int) {
	fmt.Printf("%s", "\n")
	for x := 0; x < gameWidth; x++ {
		fmt.Printf("~")
	}
	fmt.Printf("%s", "\n")
}

func printMiddleLine(gameWidth int) {
	padding := 30
	fmt.Printf("%s", "\n")
	for x := 0; x < padding; x++ {
		fmt.Printf(" ")
	}
	for x := 0; x < gameWidth-padding*2; x++ {
		fmt.Printf("_")
	}
	fmt.Printf("%s", "\n\n\n")
}

// 0 or 1 depending on who has priority
func (g *Game) PriorityIndex() int {
	for i, player := range g.Players {
		if player == g.Priority() {
			return i
		}
	}
	panic("game is corrupted")
}

// Pass makes the active player pass, whichever player has priority
func (g *Game) pass() {
	g.TakeAction(&Action{Type: Pass})
}

// passUntilPhase makes both players pass until the game is in the provided phase,
// or until the game is over.
func (g *Game) passUntilPhase(p Phase) {
	for g.Phase != p && !g.IsOver() {
		g.pass()
	}
}

// passTurn makes both players pass until it is the next turn's main phase, or until the game is over
func (g *Game) passTurn() {
	turn := g.Turn
	for g.Turn == turn && !g.IsOver() {
		g.pass()
	}
	g.passUntilPhase(Main1)
}

// playLand plays the first land it sees in the hand
func (g *Game) playLand() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsLand() {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playLand failed")
}

// putCreatureOnStack casts the first creature it sees in the hand
func (g *Game) putCreatureOnStackAndPass() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsCreature() {
			g.TakeAction(a)
			g.TakeAction(&Action{Type: PassPriority})
			return
		}
	}
	g.Print()
	panic("playCreature failed")
}

// playCreature plays the first creature it sees in the hand
func (g *Game) playCreature() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsCreature() {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playCreature failed")
}

func (g *Game) TakeActionAndResolve(action *Action) {
	g.TakeAction(action)
	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: PassPriority})
	if action.Card != nil && action.Card.EntersTheBattlefieldEffect != nil {
		g.TakeAction(&Action{Type: PassPriority})
		g.TakeAction(&Action{Type: PassPriority})
	}
}

// playAura plays the first aura it sees in the hand on its own creature
func (g *Game) playAura() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsEnchantCreature() && g.Permanent(a.Target).Owner == g.PriorityId {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playAura failed")
}

// does the first available block action
func (g *Game) doBlockAction() {
	for _, a := range g.Defender().BlockActions() {
		g.TakeAction(a)
		return
	}
	g.Print()
	panic("doBlockAction failed")
}

// playCreature plays the first creature action with Phyrexian
func (g *Game) playCreaturePhyrexian() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsCreature() && a.WithPhyrexian {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playCreaturePhyrexian failed")
}

// playInstant plays the first instant it sees in the hand
func (g *Game) playInstant() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsInstant() && a.Type == Play {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playInstant failed")
}

// playKickedInstant kicks the first kickable instant it sees in the hand
func (g *Game) playKickedInstant() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsInstant() && a.WithKicker {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playKickedInstant failed")
}

// playSorcery plays the first sorcery it sees in the hand
func (g *Game) playSorcery() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsSorcery() && a.Type == Play {
			g.TakeActionAndResolve(a)
			return
		}
	}
	g.Print()
	panic("playSorcery failed")
}

// attackWithEveryone passes priority when it's done attacking
func (g *Game) attackWithEveryone() {
	for {
		actions := g.Priority().AttackActions()
		if len(actions) == 0 {
			g.TakeAction(g.Priority().PassAction())
			return
		}
		g.TakeAction(actions[0])
	}
}

// playManaAbilityAction plays the first mana ability action
func (g *Game) playManaAbilityAction() {
	for _, a := range g.Priority().ManaActions() {
		g.TakeAction(a)
		return
	}
	g.Print()
	panic("playManaAbilityAction failed")
}

// playActivatedAbility plays the first activated ability action
func (g *Game) playActivatedAbility() {
	for _, a := range g.Priority().ActivatedAbilityActions(true, false) {
		g.TakeActionAndResolve(a)
		return
	}
	g.Print()
	panic("playActivatedAbility failed")
}

// Returns a binary representation of the game
func (g *Game) Serialize() []byte {
	gameJson, _ := json.Marshal(g)
	return gameJson
}

// Make a new game based on a serialized game
func DeserializeGame(bytes []byte) *Game {
	game := &Game{}
	json.Unmarshal(bytes, game)
	game.Players[0].game = game
	game.Players[1].game = game
	for _, perm := range game.Permanents {
		perm.game = game
	}
	return game
}

func CopyGame(g *Game) *Game {
	return DeserializeGame(g.Serialize())
}

func (g *Game) ActionStates() []*ActionState {
	actionStates := []*ActionState{}
	for _, action := range g.Actions(false) {
		actionState := &ActionState{
			Game:   g,
			Action: action,
		}
		actionStates = append(actionStates, actionState)
	}
	return actionStates
}
