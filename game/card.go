package game

import ()

// Card should be treated as immutable.
// The properties on Card are the properties like "base toughness" that do not change
// over time for a particular card.
type Card struct {
	ActivatedAbility     *Effect
	ActivatedAbilityCost *Effect
	AddsTemporaryEffect  bool
	Bloodthirst          int
	CastingCost          *CastingCost
	Effect               *Effect
	EntersPlayEffect     *Effect
	Flying               bool
	GroundEvader         bool // only blockable by fliers (like Silhana Ledgewalker)
	Hexproof             bool
	Kicker               *Effect
	Lifelink             bool
	Morbid               *Effect
	Name                 CardName
	PhyrexianCastingCost *CastingCost
	Powermenace          bool // only blockable by >= power (like Skarrgan Pitskulk)

	// http://mtg.wikia.com/wiki/Card_Types
	Subtype    []Subtype
	Supertype  []Supertype
	TargetType *TargetType
	Type       []Type

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

	BurningTreeEmissary
	EldraziSpawnToken
	Forest
	GrizzlyBears
	HungerOfTheHowlpack
	NestInvader
	NettleSentinel
	QuirionRanger
	Rancor
	SilhanaLedgewalker
	SkarrganPitskulk
	VaultSkirge
	VinesOfVastwood
)

const CARD_HEIGHT = 5
const CARD_WIDTH = 11

var Cards = map[CardName]*Card{

	/*
		Creature — Human Shaman
		When Burning-Tree Emissary enters the battlefield, add RG.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=366467
	*/
	BurningTreeEmissary: &Card{
		BasePower:        2,
		BaseToughness:    2,
		CastingCost:      &CastingCost{Colorless: 2},
		EntersPlayEffect: &Effect{Colorless: 2},
		Type:             []Type{Creature},
	},

	EldraziSpawnToken: &Card{
		BasePower:         0,
		BaseToughness:     1,
		CastingCost:       &CastingCost{Colorless: 0},
		Colorless:         1,
		SacrificesForMana: true,
		Type:              []Type{Creature},
	},

	/*
		G
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=443154
	*/
	Forest: &Card{
		Colorless: 1,
		Subtype:   []Subtype{LandForest},
		Supertype: []Supertype{Basic},
		Type:      []Type{Land},
	},

	/*
		No card text.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=4300
	*/
	GrizzlyBears: &Card{
		BasePower:     2,
		BaseToughness: 2,
		CastingCost:   &CastingCost{Colorless: 2},
		Type:          []Type{Creature},
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
		Morbid: &Effect{
			Plus1Plus1Counters: 2,
		},
		Type: []Type{Instant},
	},

	/*
		Creature — Eldrazi Drone
		When Nest Invader enters the battlefield, create a 0/1 colorless
		Eldrazi Spawn creature token. It has "Sacrifice this creature:
		Add (1)."
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=193420
	*/
	NestInvader: &Card{
		BasePower:        2,
		BaseToughness:    2,
		CastingCost:      &CastingCost{Colorless: 2},
		EntersPlayEffect: &Effect{Summon: EldraziSpawnToken},
		Type:             []Type{Creature},
	},

	/*
		Nettle Sentinel doesn't untap during your untap step.
		Whenever you cast a green spell, you may untap Nettle Sentinel.
	*/
	NettleSentinel: &Card{
		BasePower:     2,
		BaseToughness: 2,
		CastingCost:   &CastingCost{Colorless: 1},
		Type:          []Type{Creature},
	},

	/*
		Creature — Elf
		Return a Forest you control to its owner's hand: Untap target creature.
		Activate this ability only once each turn.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=3674
	*/
	QuirionRanger: &Card{
		ActivatedAbility: &Effect{
			Cost: &Effect{
				EffectType: ReturnToHand,
				TargetType: &TargetType{Subtype: LandForest, ControlledBy: SamePlayer},
			},
			EffectType: Untap,
			TargetType: &TargetType{Type: Creature},
		},
		BasePower:     1,
		BaseToughness: 1,
		CastingCost:   &CastingCost{Colorless: 1},
		Type:          []Type{Creature},
	},

	/*
		Enchanted creature gets +2/+0 and has trample.
		When Rancor is put into a graveyard from the battlefield,
		return Rancor to its owner's hand.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=442175
	*/
	Rancor: &Card{
		BasePower:     2,
		BaseToughness: 0,
		CastingCost:   &CastingCost{Colorless: 1},
		TargetType:    &TargetType{Type: Creature},
		Type:          []Type{Enchantment},
	},

	/*
		Creature — Elf Rogue
		Hexproof (This creature can't be the target of spells or abilities your
		opponents control.)
		Silhana Ledgewalker can't be blocked except by creatures with flying.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=423502
	*/
	SilhanaLedgewalker: &Card{
		BasePower:     1,
		BaseToughness: 1,
		CastingCost:   &CastingCost{Colorless: 2},
		Hexproof:      true,
		GroundEvader:  true,
		Type:          []Type{Creature},
	},

	/*
		Creature — Human Warrior
		Bloodthirst 1 (If an opponent was dealt damage this turn, this creature enters
		the battlefield with a +1/+1 counter on it.)
		Creatures with power less than Skarrgan Pit-Skulk's power can't block it.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=426622
	*/
	SkarrganPitskulk: &Card{
		BasePower:     1,
		BaseToughness: 1,
		Bloodthirst:   1,
		CastingCost:   &CastingCost{Colorless: 1},
		Powermenace:   true,
		Type:          []Type{Creature},
	},

	/*

		Artifact Creature — Imp
		(Phyrexian Black can be paid with either Black or 2 life.)
		Flying
		Lifelink (Damage dealt by this creature also causes you to gain that much life.)
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=217984
	*/
	VaultSkirge: &Card{
		BasePower:            1,
		BaseToughness:        1,
		CastingCost:          &CastingCost{Colorless: 2},
		Flying:               true,
		Hexproof:             true,
		Lifelink:             true,
		PhyrexianCastingCost: &CastingCost{Life: 2, Colorless: 1},
		Type:                 []Type{Artifact, Creature},
	},

	/*
		Kicker Green (You may pay an additional Green as you cast this spell.)
		Target creature can't be the target of spells or abilities your opponents
		control this turn. If this spell was kicked, that creature gets +4/+4 until
		end of turn.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=397747
	*/
	VinesOfVastwood: &Card{
		AddsTemporaryEffect: true,
		CastingCost:         &CastingCost{Colorless: 1},
		Kicker: &Effect{
			CastingCost: &CastingCost{Colorless: 2},
			Power:       4,
			Toughness:   4,
		},
		Effect: &Effect{
			Untargetable: true,
		},
		Type: []Type{Instant},
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

func (c *Card) IsCreature() bool {
	for _, t := range c.Type {
		if t == Creature {
			return true
		}
	}
	return false
}

func (c *Card) IsInstant() bool {
	for _, t := range c.Type {
		if t == Instant {
			return true
		}
	}
	return false
}

func (c *Card) IsLand() bool {
	for _, t := range c.Type {
		if t == Land {
			return true
		}
	}
	return false
}

func (c *Card) IsEnchantment() bool {
	for _, t := range c.Type {
		if t == Enchantment {
			return true
		}
	}
	return false
}

func (c *Card) IsEnchantCreature() bool {
	return c.IsEnchantment() && c.TargetType.Type == Creature
}

func (c *Card) HasSubtype(subtype Subtype) bool {
	for _, st := range c.Subtype {
		if st == subtype {
			return true
		}
	}
	return false
}
