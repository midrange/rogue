package game

import ()

type Player struct {
	Life  int
	Hand  []*Card
	Board []*Card

	// 0 = on the play, 1 = on the draw
	Index int
}

func (p *Player) AvailableMana() int {
	answer := 0
	for _, card := range p.Board {
		if card.IsLand {
			answer += 1
		}
	}
	return answer
}
