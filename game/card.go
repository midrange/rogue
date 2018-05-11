package game

import (
	"log"
	"math/rand"
)

type Card struct {
	// Things that are relevant wherever the card is
	Name              CardName
	IsLand            bool
	IsCreature        bool
	IsEnchantCreature bool
	ManaCost          int
	Owner             *Player

	// Properties that are relevant for any permanent
	Tapped bool
	Auras  []*Card

	// Creature-specific properties
	Attacking   bool
	Blocking    *Card
	DamageOrder []*Card
	Damage      int

	// Auras, equipment, instants, and sorceries can have targets
	Target *Card

	// For creatures these are natural.
	// For auras and equipment these indicate the boost the target gets.
	BasePower     int
	BaseToughness int
}

type CardName int

const (
	Forest CardName = iota
	GrizzlyBears
	Rancor
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
			IsCreature:    true,
			BasePower:     2,
			BaseToughness: 2,
			ManaCost:      2,
		}
	case Rancor:
		return &Card{
			IsEnchantCreature: true,
			BasePower:         2,
			BaseToughness:     0,
			ManaCost:          1,
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

func (c *Card) Power() int {
	answer := c.BasePower
	for _, aura := range c.Auras {
		answer += aura.BasePower
	}
	return answer
}

func (c *Card) Toughness() int {
	answer := c.BaseToughness
	for _, aura := range c.Auras {
		answer += aura.BaseToughness
	}
	return answer
}
