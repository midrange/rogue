package game

import (
	"fmt"
	"log"
	"math/rand"
)

type Card struct {
	Name       CardName
	IsLand     bool
	IsCreature bool
	Power      int
	Toughness  int
	ManaCost   int
	Tapped     bool
	Attacking  bool
}

//go:generate stringer -type=CardName
type CardName int

const (
	Forest CardName = iota
	GrizzlyBears
)

func NewCard(name CardName) *Card {
	card := newCardHelper(name)
	card.Name = name
	return card
}

func newCardHelper(name CardName) *Card {
	switch name {
	case Forest:
		return &Card{
			IsLand: true,
		}
	case GrizzlyBears:
		return &Card{
			IsCreature: true,
			Power:      2,
			Toughness:  2,
		}
	default:
		log.Fatalf("unimplemented card name: %d", name)
	}
	panic("control should not reach here")
}

func RandomCard() *Card {
	names := []CardName{
		Forest,
		GrizzlyBears,
	}
	index := rand.Int() % len(names)
	return NewCard(names[index])
}

func (c *Card) Print() {
	fmt.Printf("%v %v ", c.Name, c.IsLand)
}
