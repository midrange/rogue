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

func (g *Game) Actions() []*Action {
	switch g.Phase {

	case Main1:
		fallthrough
	case Main2:
		return g.Priority.PlayActions(true)

	case DeclareAttackers:
		return g.Priority.AttackActions()

	case DeclareBlockers:
		return g.Priority.BlockActions()

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

func (g *Game) HandleCombatDamage() {
	for _, card := range g.Attacker().Board {
		if card.Attacking {
			damage := card.Power()
			if damage < 0 {
				damage = 0
			}

			if len(card.DamageOrder) > 0 {
				// Deal damage to blockers
				for _, blocker := range card.DamageOrder {
					if damage == 0 {
						break
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
			} else {
				// Deal damage to the defending player
				g.Defender().Life -= damage
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

func (g *Game) NextPhase() {
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
		g.Phase = Main1
		g.Turn++
		g.Priority = g.Priority.Opponent
		g.Priority.Draw()
	}
}

func (g *Game) TakeAction(action *Action) {
	if g.IsOver() {
		panic("cannot take action when the game is over")
	}

	if action.Type == Pass {
		g.NextPhase()
		return
	}

	switch g.Phase {

	case Main1:
		fallthrough
	case Main2:
		if action.Type != Play {
			panic("expected a play or a pass during main phase")
		}
		g.Priority.Play(action.Card)

	case DeclareAttackers:
		if action.Type != Attack {
			panic("expected an attack or a pass during DeclareAttack")
		}
		action.Card.Attacking = true

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

// Pass makes the active player pass, whichever player has priority
func (g *Game) Pass() {
	g.TakeAction(&Action{Type: Pass})
}

// PassUntilPhase makes both players pass until the game is in the provided phase,
// or until the game is over.
func (g *Game) PassUntilPhase(p Phase) {
	for g.Phase != p && !g.IsOver() {
		g.Pass()
	}
}

// PassTurn makes both players pass until it is the next turn, or until the game is over
func (g *Game) PassTurn() {
	turn := g.Turn
	for g.Turn == turn && !g.IsOver() {
		g.Pass()
	}
}
