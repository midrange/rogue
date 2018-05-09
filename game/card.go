package game

import ()

type Card interface {
	Type() CardType
}

type CardType int

const (
	Land CardType = iota
	Creature
)

func RandomCard() *Card {
	panic("TODO: implement")
}
