package game

import (
	"fmt"
)

const GAME_WIDTH = 100

type PlayerId int

const (
	OnThePlay PlayerId = iota
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

	// Non-mana, non-cost effect actions go on the stack and can be responded to before they resolve
	Stack []*Action
}

//go:generate stringer -type=CardName
type Phase int

// Phase encompasses both "phases" and "steps" as per:
// https://mtg.gamepedia.com/Turn_structure
// For example, declaring attackers is a step within the combat phase in the official
// rules. Here it is just treated as another Phase.

const (
	Main1 Phase = iota
	DeclareAttackers
	DeclareBlockers
	Main2
)

func NewGame(deckToPlay *Deck, deckToDraw *Deck) *Game {
	players := [2]*Player{
		NewPlayer(deckToPlay, OnThePlay),
		NewPlayer(deckToDraw, OnTheDraw),
	}
	g := &Game{
		Players:         players,
		Phase:           Main1,
		Turn:            0,
		PriorityId:      OnThePlay,
		NextPermanentId: PermanentId(1),
		Permanents:      make(map[PermanentId]*Permanent),
		Stack:           []*Action{},
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

	if len(g.Stack) > 0 {
		actions = append(actions, g.Priority().PlayActions(false, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(false, forHuman)...)
		topAction := g.Stack[len(g.Stack)-1]
		if topAction.Owner == g.Priority() {
			actions = append(actions, &Action{
				Type: OfferToResolveNextOnStack,
			})
		} else {
			actions = append(actions, &Action{
				Type: ResolveNextOnStack,
			})
		}

		return actions
	}

	switch g.Phase {
	case Main1:
		actions = append(actions, &Action{Type: Pass})
		actions = append(actions, g.Priority().PlayActions(true, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(true, forHuman)...)
		actions = append(actions, &Action{Type: DeclareAttack})
		return append(actions, g.Priority().ManaActions()...)
	case Main2:
		actions = append(actions, &Action{Type: Pass})
		actions = append(actions, g.Priority().PlayActions(true, forHuman)...)
		actions = append(actions, g.Priority().ActivatedAbilityActions(true, forHuman)...)
		return append(actions, g.Priority().ManaActions()...)
	case DeclareAttackers:
		return append(g.Priority().AttackActions(), g.Priority().PassAction())
	case DeclareBlockers:
		return append(g.Priority().BlockActions(), g.Priority().PassAction())
	default:
		panic("unhandled phase")
	}
}

func (g *Game) AttackerId() PlayerId {
	return PlayerId(g.Turn % 2)
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
	for _, attacker := range g.Attacker().Board {
		if attacker.Attacking {
			damage := attacker.Power()
			if damage < 0 {
				damage = 0
			}

			if len(attacker.DamageOrder) > 0 {
				// Deal damage to blockers
				for _, blocker := range attacker.DamageOrder {
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

	switch g.Phase {
	case Main1:
		g.Phase = DeclareAttackers
	case DeclareAttackers:
		g.Phase = DeclareBlockers
		g.PriorityId = g.DefenderId()
	case DeclareBlockers:
		g.HandleCombatDamage()
		g.Attacker().EndCombat()
		g.Defender().EndCombat()
		g.Phase = Main2
		g.PriorityId = g.AttackerId()
	case Main2:
		// End the turn
		for _, p := range g.Players {
			p.EndTurn()
		}
		g.Phase = Main1
		g.Turn++
		g.Priority().Untap()
		g.PriorityId = g.PriorityId.OpponentId()
		g.Priority().Draw()
	}
}

func (g *Game) TakeAction(action *Action) {
	if g.IsOver() {
		panic("cannot take action when the game is over")
	}

	if action.Type == OfferToResolveNextOnStack {
		g.PriorityId = g.PriorityId.OpponentId()
		return
	} else if action.Type == ResolveNextOnStack {
		g.PriorityId = g.PriorityId.OpponentId()
		stackAction := g.Stack[len(g.Stack)-1]
		if stackAction.Type == Play {
			stackAction.Owner.ResolveSpell(stackAction)
		} else if stackAction.Type == Activate {
			stackAction.Owner.ResolveActivatedAbility(stackAction)
		}
		g.Stack = g.Stack[:len(g.Stack)-1]
		return
	}

	if action.Type == Pass {
		g.nextPhase()
		return
	}

	if action.Type == UseForMana {
		action.Source.UseForMana()
		return
	}

	switch g.Phase {

	case Main1:
		if action.Type == DeclareAttack {
			g.nextPhase()
			break
		}
		fallthrough
	case Main2:
		if action.Type == Play {
			if action.Card.IsLand() {
				g.Priority().PlayLand(action)
			} else {
				g.Stack = append(g.Stack, action)
				g.Priority().PayCostsAndPutSpellOnStack(action)
			}
		} else if action.Type == Activate {
			g.Stack = append(g.Stack, action)
			g.Priority().PayCostsAndPutAbilityOnStack(action)
		} else {
			panic("expected a play, activate, declare attack, or pass during main phase")
		}

	case DeclareAttackers:
		if action.Type != Attack {
			panic("expected an attack or a pass during DeclareAttackers")
		}
		action.With.Attacking = true
		action.With.Tapped = true

	case DeclareBlockers:
		if action.Type != Block {
			panic("expected a block or a pass during DeclareBlockers")
		}
		action.With.Blocking = action.Target
		action.Target.DamageOrder = append(action.Target.DamageOrder, action.With)

	default:
		panic("unhandled phase")
	}
}

func (g *Game) IsOver() bool {
	return g.Attacker().Lost() || g.Defender().Lost()
}

// All permanents added to the game should be created via newPermanent.
// This assigns a unique id to the permanent and activates any coming-into-play
// effects.
func (g *Game) newPermanent(card *Card, owner *Player) *Permanent {
	perm := &Permanent{
		Card:       card,
		Owner:      owner,
		TurnPlayed: g.Turn,
		Id:         g.NextPermanentId,
	}
	g.Permanents[g.NextPermanentId] = perm
	g.NextPermanentId++
	owner.Board = append(owner.Board, perm)
	perm.HandleComingIntoPlay()
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

// passTurn makes both players pass until it is the next turn, or until the game is over
func (g *Game) passTurn() {
	turn := g.Turn
	for g.Turn == turn && !g.IsOver() {
		g.pass()
	}
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

// playCreature plays the first creature it sees in the hand
func (g *Game) playCreature() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsCreature() {
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
			return
		}
	}
	g.Print()
	panic("playCreature failed")
}

// playAura plays the first aura it sees in the hand on its own creature
func (g *Game) playAura() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsEnchantCreature() && a.Target.Owner == g.Priority() {
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
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
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
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
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
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
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
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
			g.TakeAction(a)
			g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
			g.TakeAction(&Action{Type: ResolveNextOnStack})
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
		g.TakeAction(a)
		g.TakeAction(&Action{Type: OfferToResolveNextOnStack})
		g.TakeAction(&Action{Type: ResolveNextOnStack})
		return
	}
	g.Print()
	panic("playActivatedAbility failed")
}
