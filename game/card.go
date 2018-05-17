package game

import ()

// Card should be treated as immutable.
// The properties on Card are the properties like "base toughness" that do not change
// over time for a particular card.
type Card struct {
	AddsTemporaryEffect  bool
	Bloodthirst          int
	CastingCost          *CastingCost
	Effect               *Effect
	EntersPlayEffect     *Effect
	Flying               bool
	GroundEvader         bool // only blockable by fliers (like Silhana Ledgewalker)
	Hexproof             bool
	IsLand               bool
	IsCreature           bool
	IsEnchantCreature    bool
	IsInstant            bool
	Kicker               *Effect
	Lifelink             bool
	Morbid               *Effect
	Name                 CardName
	PhyrexianCastingCost *CastingCost
	Powermenace          bool // only blockable by >= power (like Skarrgan Pitskulk)

	// The base properties of creatures.
	BasePower     int
	BaseToughness int
	BaseTrample   bool

	// Properties that are relevant for Lands and other mana producers
	Colorless         int
	SacrificesForMana bool
}

//go:generate stringer -type=CardName
type CardName int

// Keep NoCard first, the rest in alphabetical order.
const (
	NoCard CardName = iota

	EldraziSpawnToken
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

var Cards = map[CardName]*Card{
	EldraziSpawnToken: &Card{
		BasePower:         0,
		BaseToughness:     1,
		CastingCost:       &CastingCost{Colorless: 0},
		Colorless:         1,
		IsCreature:        true,
		SacrificesForMana: true,
	},

	Forest: &Card{
		Colorless: 1,
		IsLand:    true,
	},

	GrizzlyBears: &Card{
		BasePower:     2,
		BaseToughness: 2,
		CastingCost:   &CastingCost{Colorless: 2},
		IsCreature:    true,
	},

	NestInvader: &Card{
		BasePower:        2,
		BaseToughness:    2,
		CastingCost:      &CastingCost{Colorless: 2},
		EntersPlayEffect: &Effect{Summon: EldraziSpawnToken},
		IsCreature:       true,
	},

	NettleSentinel: &Card{
		BasePower:     2,
		BaseToughness: 2,
		CastingCost:   &CastingCost{Colorless: 1},
		IsCreature:    true,
	},

	/*
		Enchanted creature gets +2/+0 and has trample.
		When Rancor is put into a graveyard from the battlefield,
		return Rancor to its owner's hand.
	*/
	Rancor: &Card{
		BasePower:         2,
		BaseToughness:     0,
		CastingCost:       &CastingCost{Colorless: 1},
		IsEnchantCreature: true,
	},

	/*
		Hexproof (This creature can't be the target of spells or abilities your
		opponents control.)
		Silhana Ledgewalker can't be blocked except by creatures with flying.
	*/
	SilhanaLedgewalker: &Card{
		BasePower:     1,
		BaseToughness: 1,
		CastingCost:   &CastingCost{Colorless: 2},
		Hexproof:      true,
		IsCreature:    true,
		GroundEvader:  true,
	},

	SkarrganPitskulk: &Card{
		BasePower:     1,
		BaseToughness: 1,
		Bloodthirst:   1,
		CastingCost:   &CastingCost{Colorless: 1},
		IsCreature:    true,
		Powermenace:   true,
	},

	VaultSkirge: &Card{
		BasePower:            1,
		BaseToughness:        1,
		CastingCost:          &CastingCost{Colorless: 2},
		Flying:               true,
		Hexproof:             true,
		IsCreature:           true,
		Lifelink:             true,
		PhyrexianCastingCost: &CastingCost{Life: 2, Colorless: 1},
	},

	VinesOfVastwood: &Card{
		AddsTemporaryEffect: true,
		CastingCost:         &CastingCost{Colorless: 1},
		IsInstant:           true,
		Kicker: &Effect{
			CastingCost: &CastingCost{Colorless: 2},
			Power:       4,
			Toughness:   4,
		},
		Effect: &Effect{
			Untargetable: true,
		},
	},

	/*
		Put a +1/+1 counter on target creature.
		Morbid - Put three +1/+1 counters on that creature instead if a creature died
		this turn.
	*/
	HungerOfTheHowlpack: &Card{
		Effect: &Effect{
			Plus1Plus1Counters: 1,
		},
		CastingCost: &CastingCost{Colorless: 1},
		IsInstant:   true,
		Morbid: &Effect{
			Plus1Plus1Counters: 2,
		},
	},
}

func init() {
	for name, card := range Cards {
		card.Name = name
	}
}

func (cn CardName) Card() *Card {
	return Cards[cn]
}

func (c *Card) String() string {
	p := &Permanent{Card: c}
	return p.String()
}
