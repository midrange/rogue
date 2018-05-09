package game

import ()

type Player struct {
	life  int
	hand  []*Card
	board []*Card

	// 0 = on the play, 1 = on the draw
	index int
}
