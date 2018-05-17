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
}

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
		Players:    players,
		Phase:      Main1,
		Turn:       0,
		PriorityId: OnThePlay,
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
	switch g.Phase {
	case Main1:
		actions = append(actions, g.Priority().PlayActions(true, forHuman)...)
		actions = append(actions, &Action{Type: DeclareAttack})
		return append(actions, g.Priority().ManaActions()...)
	case Main2:
		actions = g.Priority().PlayActions(true, forHuman)
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
						g.Defender().RemoveFromBoard(blocker)
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
				g.Attacker().RemoveFromBoard(attacker)
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

func (g *Game) nextPhase() {
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
			g.Priority().Play(action)
		} else {
			panic("expected a play, declare attack, or pass during main phase")
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
		if a.Card != nil && a.Card.IsLand {
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
		if a.Card != nil && a.Card.IsCreature {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playCreature failed")
}

// playCreature plays the first creature action with Phyrexian
func (g *Game) playCreaturePhyrexian() {
	for _, a := range g.Priority().PlayActions(true, false) {
		fmt.Println(a)
		if a.Card != nil && a.Card.IsCreature && a.WithPhyrexian {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playCreaturePhyrexian failed")
}

// playInstant plays the first instant it sees in the hand
func (g *Game) playInstant() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsInstant && a.Type == Play {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playInstant failed")
}

// playKickedInstant kicks the first kickable instant it sees in the hand
func (g *Game) playKickedInstant() {
	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.IsInstant && a.WithKicker {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playKickedInstant failed")
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
