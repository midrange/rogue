package game

import ()

type Game struct {
	// Players are generally referred to by index, 0 or 1.
	// Player 0 is the player who plays first.
	players [2]*Player
	phase   Phase

	// Which turn of the game is. This starts at 0.
	// In general, it is player (turn % 2)'s turn.
	turn int

	// Which player has priority
	priority int
}

type Phase int

const (
	Main1 Phase = iota
	DeclareAttack
	CombatDamage
	Main2
)

func NewGame() *Game {
	panic("TODO: implement")
}

func LegalActions() []*Action {
	panic("TODO: implement")
}

func TakeAction(action *Action) {
	panic("TODO: implement")
}
