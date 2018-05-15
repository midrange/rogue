package game

import (
	"fmt"
)

const GAME_WIDTH = 100

type Game struct {
	// Players are sometimes referred to by index, 0 or 1.
	// Player 0 is the player who plays first.
	Players [2]*Player
	Phase   Phase

	// Which turn of the game is. This starts at 0.
	// In general, it is player (turn % 2)'s turn.
	Turn int

	// Which player has priority
	Priority *Player
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
		NewPlayer(deckToPlay),
		NewPlayer(deckToDraw),
	}
	g := &Game{
		Players:  players,
		Phase:    Main1,
		Turn:     0,
		Priority: players[0],
	}

	players[0].Opponent = players[1]
	players[1].Opponent = players[0]

	players[0].Game = g
	players[1].Game = g

	return g
}

func (g *Game) Actions(forHuman bool) []*Action {
	actions := []*Action{}
	switch g.Phase {
	case Main1:
		actions = append(actions, g.Priority.PlayActions(true, forHuman)...)
		if g.canAttack() {
			actions = append(actions, &Action{Type: DeclareAttack})
		}
		return append(actions, g.Priority.ManaActions()...)
	case Main2:
		actions = g.Priority.PlayActions(true, forHuman)
		return append(actions, g.Priority.ManaActions()...)
	case DeclareAttackers:
		return append(g.Priority.AttackActions(), g.Priority.PassAction())
	case DeclareBlockers:
		return append(g.Priority.BlockActions(), g.Priority.PassAction())
	default:
		panic("unhandled phase")
	}
}

func (g *Game) Attacker() *Player {
	return g.Players[g.Turn%2]
}

func (g *Game) Defender() *Player {
	return g.Attacker().Opponent
}

func (g *Game) canAttack() bool {
	if g.Phase != Main1 {
		return false
	}
	for _, card := range g.Priority.Board {
		if card.CanAttack(g) {
			return true
		}
	}
	return false
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

func (g *Game) Creatures() []*Card {
	answer := g.Priority.Creatures()
	for _, card := range g.Priority.Opponent.Creatures() {
		answer = append(answer, card)
	}
	return answer
}

func (g *Game) nextPhase() {
	switch g.Phase {
	case Main1:
		g.Phase = DeclareAttackers
	case DeclareAttackers:
		g.Phase = DeclareBlockers
		g.Priority = g.Defender()
	case DeclareBlockers:
		g.HandleCombatDamage()
		g.Attacker().EndCombat()
		g.Defender().EndCombat()
		g.Phase = Main2
		g.Priority = g.Attacker()
	case Main2:
		// End the turn
		for _, p := range g.Players {
			p.EndTurn()
		}
		g.Phase = Main1
		g.Turn++
		g.Priority.Untap()
		g.Priority = g.Priority.Opponent
		g.Priority.Draw()
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
		action.Card.UseForMana()
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
			g.Priority.Play(action)
		} else {
			panic("expected a play, declare attack, or pass during main phase")
		}

	case DeclareAttackers:
		if action.Type != Attack {
			panic("expected an attack or a pass during DeclareAttackers")
		}
		action.Card.Attacking = true
		action.Card.Tapped = true

	case DeclareBlockers:
		if action.Type != Block {
			panic("expected a block or a pass during DeclareBlockers")
		}
		action.Card.Blocking = action.Target
		action.Target.DamageOrder = append(action.Target.DamageOrder, action.Card)

	default:
		panic("unhandled phase")
	}
}

func (g *Game) IsOver() bool {
	return g.Priority.Lost() || g.Priority.Opponent.Lost()
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
	fmt.Printf("%v", "\n")
	for x := 0; x < gameWidth; x++ {
		fmt.Printf("~")
	}
	fmt.Printf("%v", "\n")
}

func printMiddleLine(gameWidth int) {
	padding := 30
	fmt.Printf("%v", "\n")
	for x := 0; x < padding; x++ {
		fmt.Printf(" ")
	}
	for x := 0; x < gameWidth-padding*2; x++ {
		fmt.Printf("_")
	}
	fmt.Printf("%v", "\n\n\n")
}

// 0 or 1 depending on who has priority
func (g *Game) PriorityIndex() int {
	for i, player := range g.Players {
		if player == g.Priority {
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
	for _, a := range g.Priority.PlayActions(true, false) {
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
	for _, a := range g.Priority.PlayActions(true, false) {
		if a.Card != nil && a.Card.IsCreature {
			g.TakeAction(a)
			return
		}
	}
	g.Print()
	panic("playCreature failed")
}

// playInstant plays the first instant it sees in the hand
func (g *Game) playInstant() {
	for _, a := range g.Priority.PlayActions(true, false) {
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
	for _, a := range g.Priority.PlayActions(true, false) {
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
		actions := g.Priority.AttackActions()
		if len(actions) == 0 {
			g.TakeAction(g.Priority.PassAction())
			return
		}
		g.TakeAction(actions[0])
	}
}
