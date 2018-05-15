package game

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

type Card struct {
	// Things that are relevant wherever the card is
	AddsTemporaryEffect  bool
	CastingCost          *CastingCost
	EntersPlayAction     *Action
	HasEntersPlayAction  bool
	HasKicker            bool
	HasManaAbility       bool
	HasMorbid            bool
	HasPhyrexian         bool
	IsLand               bool
	IsCreature           bool
	IsEnchantCreature    bool
	IsInstant            bool
	IsToken              bool
	Kicker               *Modifier
	Modifier             *Modifier // temporary Effect Modifiers for Power, Toughness, Untargetable, Hexproof
	Morbid               *Modifier
	Name                 CardName
	Owner                *Player
	PhyrexianCastingCost *CastingCost

	// Properties that are relevant for any permanent
	Auras      []*Card
	Effects    []*Effect
	Tapped     bool
	TurnPlayed int

	// Properties that are relevant for Lands and other mana producers
	Colorless         int
	SacrificesForMana bool

	// Creature-specific properties
	Attacking         bool
	Blocking          *Card
	Bloodthirst       int
	DamageOrder       []*Card
	Damage            int
	Flying            bool
	GroundEvader      bool // a fake keyword for Silhana Ledgewalker that faux flies
	Hexproof          bool
	Lifelink          bool
	PowerCounters     int
	Powermenace       bool // a fake keyword for Skarrgan Pitskulk, who can only be blocked by >= Power
	ToughnessCounters int

	// Auras, equipment, instants, and sorceries can have targets
	Target *Card

	// For creatures these are natural.
	// For auras and equipment, these indicate the boost the target gets.
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
			IsToken:           true,
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
			CastingCost: &CastingCost{Colorless: 1},
			HasMorbid:   true,
			IsInstant:   true,
			Morbid: &Modifier{
				PowerCounters:     2,
				ToughnessCounters: 2,
			},
			PowerCounters:     1,
			ToughnessCounters: 1,
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

func (c *Card) String() string {
	if c.IsLand {
		return fmt.Sprintf("%v", c.Name)
	} else if c.IsCreature {
		return fmt.Sprintf("%v (%v/%v)", c.Name, c.Power(), c.Toughness())
	}
	return fmt.Sprintf("%v", c.Name)
}

func (c *Card) AsciiImage(showBack bool) [CARD_HEIGHT][CARD_WIDTH]string {
	const cardWidth = CARD_WIDTH
	const cardHeight = CARD_HEIGHT
	imageGrid := [cardHeight][cardWidth]string{}
	for y := 0; y < cardHeight; y++ {
		for x := 0; x < cardWidth; x++ {
			if y == 0 || y == cardHeight-1 {
				imageGrid[y][x] = string('-')
			} else if x == 0 || x == cardWidth-1 {
				imageGrid[y][x] = string('|')
			} else {
				imageGrid[y][x] = string(' ')
			}
		}
	}

	initialIndex := 2

	if showBack {
		middleX := cardWidth / 2
		middleY := cardHeight / 2

		noon := []int{middleX, middleY - 1}
		two := []int{middleX + 2, middleY}
		ten := []int{middleX - 2, middleY}
		seven := []int{middleX - 1, middleY + 1}
		four := []int{middleX + 1, middleY + 1}

		points := [][]int{noon, two, four, seven, ten}
		for _, p := range points {
			imageGrid[p[1]][p[0]] = string('*')
		}
	} else {
		nameRow := 2
		words := strings.Split(fmt.Sprintf("%v", c.Name), " ")
		for _, word := range words {
			wordWidth := Min(3, len(word))
			if len(words) == 1 {
				wordWidth = Min(len(word), cardWidth-4)
			}
			for x := initialIndex; x < wordWidth+initialIndex; x++ {
				imageGrid[nameRow][x] = string(word[x-initialIndex])
			}
			initialIndex += wordWidth + 1
			if initialIndex >= cardWidth-wordWidth-1 {
				break
			}
		}

		if c.IsCreature {
			initialIndex := 2
			statsRow := 3
			statsString := fmt.Sprintf("%v/%v", c.Power(), c.Toughness())
			for x := initialIndex; x < len(statsString)+initialIndex; x++ {
				imageGrid[statsRow][x] = string(statsString[x-initialIndex])
			}

		}

		if !c.IsLand {
			initialIndex := 2
			ccRow := 1
			ccString := fmt.Sprintf("%v", c.CastingCost.Colorless)
			for x := initialIndex; x < len(ccString)+initialIndex; x++ {
				imageGrid[ccRow][x] = string(ccString[x-initialIndex])
			}
		}

		if c.Tapped {
			tappedRow := 0
			initialIndex := 0
			tappedString := "TAPPED"
			for x := initialIndex; x < len(tappedString)+initialIndex; x++ {
				imageGrid[tappedRow][x] = string(tappedString[x-initialIndex])
			}
		}
	}

	return imageGrid
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (c *Card) Power() int {
	answer := c.BasePower + c.PowerCounters
	for _, aura := range c.Auras {
		answer += aura.BasePower
	}
	for _, effect := range c.Effects {
		answer += effect.Card.Modifier.Power
		if effect.Action.WithKicker {
			answer += effect.Card.Kicker.Power
		}
	}
	return answer
}

func (c *Card) Toughness() int {
	answer := c.BaseToughness + c.ToughnessCounters
	for _, aura := range c.Auras {
		answer += aura.BaseToughness
	}
	for _, effect := range c.Effects {
		answer += effect.Card.Modifier.Toughness
		if effect.Action.WithKicker {
			answer += effect.Card.Kicker.Toughness
		}
	}
	return answer
}

func (c *Card) Targetable(targetingSpell *Card) bool {
	answer := true
	if targetingSpell.Owner != c.Owner && c.Hexproof {
		return false
	}
	for _, effect := range c.Effects {
		if effect.Card.Modifier.Untargetable == true {
			return false
		}
		if targetingSpell.Owner != c.Owner && effect.Card.Modifier.Hexproof == true {
			return false
		}
	}
	return answer
}

func (c *Card) CanAttack(g *Game) bool {
	if c.Tapped || !c.IsCreature || c.Power() == 0 || c.TurnPlayed == g.Turn {
		return false
	}
	return true
}

func (c *Card) Trample() bool {
	if c.BaseTrample {
		return true
	}
	for _, aura := range c.Auras {
		if aura.BaseTrample {
			return true
		}
	}
	return false
}

func (c *Card) RespondToUntapPhase() {
	if c.Name != NettleSentinel {
		c.Tapped = false
	}
}

func (c *Card) RespondToSpell(spell *Card) {
	if c.Name == NettleSentinel {
		c.Tapped = false
	}
}

func (c *Card) ManaActions() []*Action {
	// TODO handle lands you can also sac for mana
	if (c.IsLand && !c.Tapped) || c.HasManaAbility {
		return []*Action{&Action{Type: UseForMana, Card: c}}
	}
	return []*Action{}
}

func (c *Card) UseForMana() {
	c.Owner.AddMana(c.Colorless)
	c.Tapped = true
	if c.SacrificesForMana {
		c.Owner.RemoveFromBoard(c)
	}
}

func (c *Card) DoEffect(action *Action) {
	if c.AddsTemporaryEffect {
		action.Target.Effects = append(action.Target.Effects, &Effect{Action: action, Card: c})
	}
	if action.Target != nil {
		// note that Counters and Morbid Counters are additive
		action.Target.PowerCounters += c.PowerCounters
		action.Target.ToughnessCounters += c.ToughnessCounters
		if c.HasMorbid && (c.Owner.CreatureDied || c.Owner.Opponent.CreatureDied) {
			action.Target.PowerCounters += c.Morbid.PowerCounters
			action.Target.ToughnessCounters += c.Morbid.ToughnessCounters
		}
	}
}

func (c *Card) HasLegalTarget(g *Game) bool {
	for _, creature := range g.Creatures() {
		if creature.Targetable(c) {
			return true
		}
	}
	return false
}

func (c *Card) CanBlock(attacker *Card) bool {
	if attacker.GroundEvader && !c.Flying {
		return false
	}
	if attacker.Flying && !c.Flying {
		return false
	}
	if attacker.Powermenace && attacker.Power() > c.Power() {
		return false
	}
	return true
}

func (c *Card) DoComesIntoPlayEffects() {
	if c.Bloodthirst > 0 && c.Owner.Opponent.DamageThisTurn > 0 {
		c.PowerCounters += c.Bloodthirst
		c.ToughnessCounters += c.Bloodthirst
	}
	if c.HasEntersPlayAction {
		c.Owner.Game.TakeAction(c.EntersPlayAction)
	}
}

/*
	Most creatures don't do anything special when they deal damage.
	Currently just ones with Lifelink do something extra.
*/
func (c *Card) DidDealDamage(damage int) {
	if c.Lifelink && damage > 0 {
		c.Owner.Life += damage
	}
}
