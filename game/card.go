package game

import (
	"log"
)

// Card should be treated as immutable.
// The properties on Card are the properties like "base toughness" that do not change
// over time for a particular card.
type Card struct {
	AddsTemporaryEffect  bool
	Bloodthirst          int
	CastingCost          *CastingCost
	EntersPlayAction     *Action
	HasEntersPlayAction  bool
	HasKicker            bool
	HasManaAbility       bool
	Flying               bool
	GroundEvader         bool // only blockable by fliers (like Silhana Ledgewalker)
	Hexproof             bool
	HasMorbid            bool
	HasPhyrexian         bool
	HasKicker            bool
	IsLand               bool
	IsCreature           bool
	IsEnchantCreature    bool
	IsInstant            bool
	Kicker               *Modifier
	Lifelink             bool
	ManaCost             int
	Modifier             *Modifier
	Morbid               *Modifier
	Name                 CardName
	PhyrexianCastingCost *CastingCost
	Powermenace          bool // only blockable by >= power (like Skarrgan Pitskulk)

	// The base properties of creatures.
	BasePower     int
	BaseToughness int
	BaseTrample   bool
}

//go:generate stringer -type=CardName
type CardName int

const (
	EldraziSpawnToken CardName = iota
	Forest
	GrizzlyBears
	HungerOfTheHowlpack
	NestInvader
	NettleSentinel
	Rancor
	SilhanaLedgewalker
	SkarrganPitskulk
	VaultSkirge
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
			Colorless: 1,
			IsLand:    true,
		}
	case GrizzlyBears:
		return &Card{
			BasePower:     2,
			BaseToughness: 2,
			CastingCost:   &CastingCost{Colorless: 2},
			IsCreature:    true,
		}
	case NestInvader:
		tokenCard := &Card{
			BasePower:         0,
			BaseToughness:     1,
			CastingCost:       &CastingCost{Colorless: 0},
			Colorless:         1,
			HasManaAbility:    true,
			IsCreature:        true,
			Name:              EldraziSpawnToken,
			SacrificesForMana: true,
		}
		return &Card{
			BasePower:           2,
			BaseToughness:       2,
			CastingCost:         &CastingCost{Colorless: 2},
			EntersPlayAction:    &Action{Type: Play, Card: tokenCard},
			HasEntersPlayAction: true,
			IsCreature:          true,
		}
	case NettleSentinel:
		/*
			Nettle Sentinel doesn't untap during your untap step.
			Whenever you cast a green spell, you may untap Nettle Sentinel.
		*/
		return &Card{
			BasePower:     2,
			BaseToughness: 2,
			CastingCost:   &CastingCost{Colorless: 1},
			IsCreature:    true,
		}
	case Rancor:
		/*
			Enchanted creature gets +2/+0 and has trample.
			When Rancor is put into a graveyard from the battlefield,
			return Rancor to its owner's hand.
		*/
		return &Card{
			BasePower:         2,
			BaseToughness:     0,
			CastingCost:       &CastingCost{Colorless: 1},
			IsEnchantCreature: true,
		}
	case SilhanaLedgewalker:
		/*
			Hexproof (This creature can't be the target of spells or abilities your opponents control.)
			Silhana Ledgewalker can't be blocked except by creatures with flying.
		*/
		return &Card{
			BasePower:     1,
			BaseToughness: 1,
			CastingCost:   &CastingCost{Colorless: 2},
			Hexproof:      true,
			IsCreature:    true,
			GroundEvader:  true,
		}
	case SkarrganPitskulk:
		/*
			Bloodthirst 1 (If an opponent was dealt damage this turn, this creature enters the
			battlefield with a +1/+1 counter on it.)
			Creatures with power less than Skarrgan Pit-Skulk's power can't block it.
		*/
		return &Card{
			BasePower:     1,
			BaseToughness: 1,
			Bloodthirst:   1,
			CastingCost:   &CastingCost{Colorless: 1},
			IsCreature:    true,
			Powermenace:   true,
		}
	case VaultSkirge:
		/*
			(Phyrexian Black can be paid with either Black or 2 life.)
			Flying
			Lifelink (Damage dealt by this creature also causes you to gain that much life.)
		*/
		return &Card{
			BasePower:            1,
			BaseToughness:        1,
			CastingCost:          &CastingCost{Colorless: 2},
			Flying:               true,
			HasPhyrexian:         true,
			Hexproof:             true,
			IsCreature:           true,
			Lifelink:             true,
			PhyrexianCastingCost: &CastingCost{Life: 2, Colorless: 1},
		}
	case VinesOfVastwood:
		return &Card{
			AddsTemporaryEffect: true,
			CastingCost:         &CastingCost{Colorless: 1},
			HasKicker:           true,
			IsInstant:           true,
			Kicker: &Modifier{
				CastingCost: &CastingCost{Colorless: 2},
				Power:       4,
				Toughness:   4,
			},
			Modifier: &Modifier{
				Untargetable: true,
			},
		}
	case HungerOfTheHowlpack:
		/*
			Put a +1/+1 counter on target creature.
			Morbid - Put three +1/+1 counters on that creature instead if a creature died this turn.
		*/
		return &Card{
			Modifier: &Modifier{
				Plus1Plus1Counters: 1,
			},
			CastingCost: &CastingCost{Colorless: 1},
			HasMorbid:   true,
			IsInstant:   true,
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
