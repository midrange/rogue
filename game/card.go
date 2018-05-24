package game

import ()

// Card should be treated as immutable.
// The properties on Card are the properties like "base toughness" that do not change
// over time for a particular card.
type Card struct {
	ActivatedAbility           *Effect
	AddsTemporaryEffect        bool
	AlternateCastingCost       *Cost
	Bloodthirst                int
	CastingCost                *Cost
	Effects                    []*Effect
	EntersGraveyardEffect      *Effect
	EntersTheBattlefieldEffect *Effect
	Flash                      bool
	Flying                     bool
	GroundEvader               bool // only blockable by fliers (like Silhana Ledgewalker)
	Hexproof                   bool
	Kicker                     *Effect
	Lifelink                   bool
	Morbid                     *Effect
	Name                       CardName
	PhyrexianCastingCost       *Cost
	Powermenace                bool // only blockable by >= power (like Skarrgan Pitskulk)

	// http://mtg.wikia.com/wiki/Card_Types
	Subtype   []Subtype
	Supertype []Supertype
	Selector  *Selector
	Type      []Type

	// The base properties of creatures.
	BasePower     int
	BaseToughness int
	BaseTrample   bool

	// Properties that are relevant for Lands and other mana producers
	Colorless         int
	SacrificesForMana bool

	// Properties that are relevant for Auras
	EnchantedPermanentDiesEffect *Effect
}

//go:generate stringer -type=CardName
type CardName int

// Keep NoCard first, the rest in alphabetical order.
const (
	NoCard CardName = iota

	BurningTreeEmissary
	Counterspell
	Daze
	EldraziSpawnToken
	ElephantGuide
	ElephantToken
	FaerieMiscreant
	Forest
	GrizzlyBears
	Gush
	HungerOfTheHowlpack
	Island
	MutagenicGrowth
	NestInvader
	NettleSentinel
	QuirionRanger
	Rancor
	SilhanaLedgewalker
	SkarrganPitskulk
	Snap
	SpellstutterSprite
	VaultSkirge
	VinesOfVastwood
)

var Cards = map[CardName]*Card{

	/*
		Creature — Human Shaman
		When Burning-Tree Emissary enters the battlefield, add RG.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=366467
	*/
	BurningTreeEmissary: &Card{
		BasePower:                  2,
		BaseToughness:              2,
		CastingCost:                &Cost{Colorless: 2},
		EntersTheBattlefieldEffect: &Effect{Colorless: 2, EffectType: AddMana},
		Type: []Type{Creature},
	},

	/*
		Counter target spell.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=202437
	*/
	Counterspell: &Card{
		CastingCost: &Cost{Colorless: 2},
		Effects: []*Effect{
			&Effect{
				EffectType: Countermagic,
			},
		},
		Selector: &Selector{Type: Spell},
		Type:     []Type{Instant},
	},

	/*
		You may return an Island you control to its owner's hand rather than pay
		this spell's mana cost.

		Counter target spell unless its controller pays 1.

		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=21284
	*/
	Daze: &Card{
		AlternateCastingCost: &Cost{
			Effect: &Effect{
				EffectType: ReturnToHand,
				Selector:   &Selector{Subtype: LandIsland, ControlledBy: SamePlayer, Count: 1, Targetted: false},
			},
		},
		CastingCost: &Cost{Colorless: 2},
		Effects: []*Effect{&Effect{
			EffectType: ManaSink,
			Selector:   &Selector{Type: Spell, Count: 1},
		}},
		Type: []Type{Instant},
	},

	/*
		Created by NestInvader.
	*/
	EldraziSpawnToken: &Card{
		BasePower:         0,
		BaseToughness:     1,
		CastingCost:       &Cost{Colorless: 0},
		Colorless:         1,
		SacrificesForMana: true,
		Type:              []Type{Creature},
	},

	/*

		Enchanted creature gets +3/+3.
		When enchanted creature dies, create a 3/3 green Elephant creature token.
	*/
	ElephantGuide: &Card{
		BasePower:                    3,
		BaseToughness:                3,
		CastingCost:                  &Cost{Colorless: 3},
		EnchantedPermanentDiesEffect: &Effect{Summon: ElephantToken},
		Selector:                     &Selector{Type: Creature},
		Type:                         []Type{Enchantment},
	},

	/*
		Created by NestInvader.
	*/
	ElephantToken: &Card{
		BasePower:     3,
		BaseToughness: 3,
		CastingCost:   &Cost{Colorless: 0},
		Type:          []Type{Creature},
	},

	/*
		Flying (This creature can't be blocked except by creatures with flying or reach.)
		When Faerie Miscreant enters the battlefield, if you control another creature
		named Faerie Miscreant, draw a card.
	*/
	FaerieMiscreant: &Card{
		BasePower:                  1,
		BaseToughness:              1,
		CastingCost:                &Cost{Colorless: 1},
		EntersTheBattlefieldEffect: &Effect{Condition: &Condition{ControlAnother: FaerieMiscreant}, EffectType: DrawCard},
		Flying:  true,
		Subtype: []Subtype{Faerie},
		Type:    []Type{Creature},
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
		CastingCost:   &Cost{Colorless: 2},
		Type:          []Type{Creature},
	},

	/*
		You may return two Islands you control to their owner's hand rather than pay this spell's mana cost.
		Draw two cards.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=20404
	*/
	Gush: &Card{
		AlternateCastingCost: &Cost{
			Effect: &Effect{
				EffectType: ReturnToHand,
				Selector:   &Selector{Subtype: LandIsland, ControlledBy: SamePlayer, Count: 2, Targetted: false},
			},
		},
		CastingCost: &Cost{Colorless: 5},
		Effects: []*Effect{&Effect{
			EffectType: DrawCard,
			Selector:   &Selector{Count: 2},
		}},
		Type: []Type{Instant},
	},

	/*
		Put a +1/+1 counter on target creature.
		Morbid - Put three +1/+1 counters on that creature instead if a creature died
		this turn.
	*/
	HungerOfTheHowlpack: &Card{
		Effects: []*Effect{&Effect{
			Plus1Plus1Counters: 1,
		}},
		CastingCost: &Cost{Colorless: 1},
		Morbid: &Effect{
			Plus1Plus1Counters: 2,
		},
		Type: []Type{Instant},
	},

	/*
		U
		http://gatherer.wizards.com/Pages/Card/Details.aspx?name=ISLAND
	*/
	Island: &Card{
		Colorless: 1,
		Subtype:   []Subtype{LandIsland},
		Supertype: []Supertype{Basic},
		Type:      []Type{Land},
	},

	/*
		(Phyrexian Green can be paid with either Green or 2 life.)
		Target creature gets +2/+2 until end of turn.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?name=mutagenic%20growth
	*/
	MutagenicGrowth: &Card{
		AddsTemporaryEffect: true,
		CastingCost:         &Cost{Colorless: 1},
		Effects: []*Effect{&Effect{
			Power:     2,
			Toughness: 2,
		}},
		PhyrexianCastingCost: &Cost{Life: 2},
		Type:                 []Type{Instant},
	},

	/*
		Creature — Eldrazi Drone
		When Nest Invader enters the battlefield, create a 0/1 colorless
		Eldrazi Spawn creature token. It has "Sacrifice this creature:
		Add (1)."
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=193420
	*/
	NestInvader: &Card{
		BasePower:                  2,
		BaseToughness:              2,
		CastingCost:                &Cost{Colorless: 2},
		EntersTheBattlefieldEffect: &Effect{Summon: EldraziSpawnToken},
		Type: []Type{Creature},
	},

	/*
		Nettle Sentinel doesn't untap during your untap step.
		Whenever you cast a green spell, you may untap Nettle Sentinel.
	*/
	NettleSentinel: &Card{
		BasePower:     2,
		BaseToughness: 2,
		CastingCost:   &Cost{Colorless: 1},
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
			Cost: &Cost{
				Effect: &Effect{
					EffectType: ReturnToHand,
					Selector:   &Selector{Subtype: LandForest, ControlledBy: SamePlayer, Targetted: false},
				},
			},
			EffectType: Untap,
			Selector:   &Selector{Type: Creature},
		},
		BasePower:     1,
		BaseToughness: 1,
		CastingCost:   &Cost{Colorless: 1},
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
		CastingCost:   &Cost{Colorless: 1},
		EntersGraveyardEffect: &Effect{
			EffectType: ReturnToHand,
		},
		Selector: &Selector{Type: Creature},
		Type:     []Type{Enchantment},
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
		CastingCost:   &Cost{Colorless: 2},
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
		CastingCost:   &Cost{Colorless: 1},
		Powermenace:   true,
		Type:          []Type{Creature},
	},

	/*
		Return target creature to its owner's hand. Untap up to two lands.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?name=snap
	*/
	Snap: &Card{
		CastingCost: &Cost{Colorless: 2},
		Effects: []*Effect{
			&Effect{
				EffectType: Untap,
				Selector:   &Selector{Type: Land, Count: 2, Targetted: false},
			},
			&Effect{
				EffectType: ReturnToHand,
				Selector:   &Selector{Type: Creature, Targetted: true},
			},
		},
		Type: []Type{Instant},
	},

	/*
		Flash
		Flying
		When Spellstutter Sprite enters the battlefield, counter target spell with converted mana
		cost X or less, where X is the number of Faeries you control.
		http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=139429
	*/
	SpellstutterSprite: &Card{
		BasePower:     1,
		BaseToughness: 1,
		CastingCost:   &Cost{Colorless: 2},
		EntersTheBattlefieldEffect: &Effect{
			EffectType: Countermagic,
			Condition:  &Condition{ConvertedManaCostLTE: Faerie},
			Selector:   &Selector{Type: Spell},
		},
		Flash:   true,
		Flying:  true,
		Subtype: []Subtype{Faerie},
		Type:    []Type{Creature},
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
		CastingCost:          &Cost{Colorless: 2},
		Flying:               true,
		Lifelink:             true,
		PhyrexianCastingCost: &Cost{Life: 2, Colorless: 1},
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
		CastingCost:         &Cost{Colorless: 1},
		Kicker: &Effect{
			Cost:      &Cost{Colorless: 2},
			Power:     4,
			Toughness: 4,
		},
		Effects: []*Effect{&Effect{
			Untargetable: true,
		}},
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
	return c.IsEnchantment() && c.Selector.Type == Creature
}

func (c *Card) IsSorcery() bool {
	for _, t := range c.Type {
		if t == Sorcery {
			return true
		}
	}
	return false
}

func (c *Card) IsSpell() bool {
	return c.IsSorcery() || c.IsInstant()
}

func (c *Card) HasSubtype(subtype Subtype) bool {
	for _, st := range c.Subtype {
		if st == subtype {
			return true
		}
	}
	return false
}

func (c *Card) HasCreatureTargets() bool {
	for _, e := range c.Effects {
		if e.Selector != nil {
			if e.Selector.Type == Creature {
				return true
			}
		}
	}
	return false
}

func (c *Card) HasSpellTargets() bool {
	for _, e := range c.Effects {
		if e.Selector != nil {
			if e.Selector.Type == Spell {
				return true
			}
		}
	}
	return false
}

func (c *Card) HasEntersTheBattlefieldTargets() bool {
	if c.EntersTheBattlefieldEffect != nil {
		if c.EntersTheBattlefieldEffect.Selector != nil {
			return true
		}
	}
	return false
}
