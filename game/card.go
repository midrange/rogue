package game

import (
	"log"
)

// Card should be treated as immutable.
// The properties on Card are the properties like "base toughness" that do not change
// over time for a particular card.
type Card struct {
	AddsTemporaryEffect bool
	Bloodthirst         int
	Flying              bool
	GroundEvader        bool // only blockable by fliers (like Silhana Ledgewalker)
	Hexproof            bool
	HasMorbid           bool
	HasKicker           bool
	IsLand              bool
	IsCreature          bool
	IsEnchantCreature   bool
	IsInstant           bool
	Kicker              *Modifier
	ManaCost            int
	Modifier            *Modifier
	Morbid              *Modifier
	Name                CardName
	Powermenace         bool // only blockable by >= power (like Skarrgan Pitskulk)

	// The base properties of creatures.
	BasePower     int
	BaseToughness int
	BaseTrample   bool
}

//go:generate stringer -type=CardName
type CardName int

const (
	Forest CardName = iota
	GrizzlyBears
	HungerOfTheHowlpack
	NettleSentinel
	Rancor
	SilhanaLedgewalker
	SkarrganPitskulk
	VinesOfVastwood
)

const CARD_HEIGHT = 5
const CARD_WIDTH = 11

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
	case NettleSentinel:
		/*
			Nettle Sentinel doesn't untap during your untap step.
			Whenever you cast a green spell, you may untap Nettle Sentinel.
		*/
		return &Card{
			IsCreature:    true,
			BasePower:     2,
			BaseToughness: 2,
			ManaCost:      1,
		}
	case SilhanaLedgewalker:
		/*
			Hexproof (This creature can't be the target of spells or abilities your opponents control.)
			Silhana Ledgewalker can't be blocked except by creatures with flying.
		*/
		return &Card{
			IsCreature:    true,
			BasePower:     1,
			BaseToughness: 1,
			ManaCost:      2,
			Hexproof:      true,
			GroundEvader:  true,
		}
	case SkarrganPitskulk:
		/*
			Bloodthirst 1 (If an opponent was dealt damage this turn, this creature enters the
			battlefield with a +1/+1 counter on it.)
			Creatures with power less than Skarrgan Pit-Skulk's power can't block it.
		*/
		return &Card{
			IsCreature:    true,
			BasePower:     1,
			BaseToughness: 1,
			ManaCost:      1,
			Bloodthirst:   1,
			Powermenace:   true,
		}
	case Rancor:
		/*
			Enchanted creature gets +2/+0 and has trample.
			When Rancor is put into a graveyard from the battlefield,
			return Rancor to its owner's hand.
		*/
		return &Card{
			IsEnchantCreature: true,
			BasePower:         2,
			BaseToughness:     0,
			ManaCost:          1,
		}
	case VinesOfVastwood:
		return &Card{
			AddsTemporaryEffect: true,
			HasKicker:           true,
			IsInstant:           true,
			ManaCost:            1,
			Modifier: &Modifier{
				Untargetable: true,
			},
			Kicker: &Modifier{
				Cost:      2,
				Power:     4,
				Toughness: 4,
			},
		}
	case HungerOfTheHowlpack:
		/*
			Put a +1/+1 counter on target creature.
			Morbid - Put three +1/+1 counters on that creature instead if a creature died this turn.
		*/
		return &Card{
			HasMorbid: true,
			IsInstant: true,
			Modifier: &Modifier{
				Plus1Plus1Counters: 1,
			},
			ManaCost: 1,
			Morbid: &Modifier{
				Plus1Plus1Counters: 2,
			},
		}
	default:
		log.Fatalf("unimplemented card name: %d", name)
	}
	panic("control should not reach here")
}

func (c *Card) String() string {
	p := &Permanent{Card: c}
	return p.String()
}
