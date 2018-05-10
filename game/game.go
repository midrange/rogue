package game

import (
	"fmt"
)

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

const (
	Main1 Phase = iota
	DeclareAttack
	Main2
)

func NewGame() *Game {
	players := [2]*Player{NewPlayer(), NewPlayer()}
	players[0].Opponent = players[1]
	players[1].Opponent = players[0]
	return &Game{
		Players:  players,
		Phase:    Main1,
		Turn:     0,
		Priority: players[0],
	}
}

func (g *Game) Actions() []*Action {
	switch g.Phase {

	case Main1:
		fallthrough
	case Main2:
		return g.Priority.PlayActions(true)

	case DeclareAttack:
		return g.Priority.AttackActions()

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
			g.Defender().Life -= card.Power
		}
	}
}

func (g *Game) NextPhase() {
	switch g.Phase {
	case Main1:
		g.Phase = DeclareAttack
	case DeclareAttack:
		g.HandleCombatDamage()
		g.Attacker().EndCombat()
		g.Phase = Main2
	case Main2:
		// End the turn
		g.Phase = Main1
		g.Turn++
		g.Priority = g.Priority.Opponent
	}
}

func (g *Game) TakeAction(action *Action) {
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

	case DeclareAttack:
		if action.Type != Attack {
			panic("expected an attack or a pass during DeclareAttack")
		}
		action.Card.Attacking = true

	default:
		panic("unhandled phase")
	}
}

func (g *Game) IsOver() bool {
	for _, player := range g.Players {
		if player.Life <= 0 {
			return true
		}
	}
	return false
}

func (g *Game) Print() {
	fmt.Println("\n~~~~~~")
	g.Players[1].Print(1, false)
	g.Players[0].Print(0, false)
	fmt.Println("~~~~~~\n")	
}
