package game

import (
	"fmt"
	"strings"
)

// A PermanentId is provided for each permanent when it enters the battlefield.
// Each id is unique for a particular game and is never reused for subsequent
// permanents.
// The first allocated id is 1. This way, 0 is not a valid PermanentId, so if you
// see anything with an id of 0 it means we are using something uninitialized.
type PermanentId int

type Permanent struct {
	*Card
	Id PermanentId

	// Properties that are relevant for any permanent
	ActivatedThisTurn bool
	Auras             []*Permanent
	Effects           []*Effect
	Owner             *Player
	Tapped            bool
	TurnPlayed        int

	// Creature-specific properties
	Attacking          bool
	Blocking           *Permanent
	DamageOrder        []*Permanent
	Damage             int
	Plus1Plus1Counters int

	// Auras, equipment, instants, and sorceries can have targets
	Target *Permanent
}

func (p *Permanent) String() string {
	if p.IsLand() {
		return fmt.Sprintf("%s", p.Name)
	} else if p.IsCreature() {
		return fmt.Sprintf("%s (%d/%d)", p.Name, p.Power(), p.Toughness())
	}
	return fmt.Sprintf("%s", p.Name)
}

func (c *Permanent) AsciiImage(showBack bool) [CARD_HEIGHT][CARD_WIDTH]string {
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
		words := strings.Split(fmt.Sprintf("%s", c.Name), " ")
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

		if c.IsCreature() {
			initialIndex := 2
			statsRow := 3
			statsString := fmt.Sprintf("%d/%d", c.Power(), c.Toughness())
			for x := initialIndex; x < len(statsString)+initialIndex; x++ {
				imageGrid[statsRow][x] = string(statsString[x-initialIndex])
			}

		}

		if !c.IsLand() {
			initialIndex := 2
			ccRow := 1
			ccString := fmt.Sprintf("%d", c.CastingCost.Colorless)
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

func (p *Permanent) Power() int {
	answer := p.BasePower + p.Plus1Plus1Counters
	for _, aura := range p.Auras {
		answer += aura.BasePower
	}
	for _, effect := range p.Effects {
		answer += effect.Power
		if effect.Kicker != nil {
			answer += effect.Kicker.Power
		}
	}
	return answer
}

func (c *Permanent) Toughness() int {
	answer := c.BaseToughness + c.Plus1Plus1Counters
	for _, aura := range c.Auras {
		answer += aura.BaseToughness
	}
	for _, effect := range c.Effects {
		answer += effect.Toughness
		if effect.Kicker != nil {
			answer += effect.Kicker.Toughness
		}
	}
	return answer
}

func (c *Permanent) CanAttack(g *Game) bool {
	if c.Tapped || !c.IsCreature() || c.Power() == 0 || c.TurnPlayed == g.Turn {
		return false
	}
	return true
}

func (c *Permanent) Trample() bool {
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

func (c *Permanent) RespondToUntapPhase() {
	if c.Name != NettleSentinel {
		c.Tapped = false
	}
}

func (c *Permanent) RespondToSpell() {
	if c.Name == NettleSentinel {
		c.Tapped = false
	}
}

func (c *Permanent) ManaActions() []*Action {
	if c.IsLand() && !c.Tapped || c.SacrificesForMana {
		return []*Action{&Action{Type: UseForMana, Source: c}}
	}
	return []*Action{}
}

func (c *Permanent) UseForMana() {
	c.Owner.AddMana(c.Colorless)
	c.Tapped = true
	if c.SacrificesForMana {
		c.Owner.RemoveFromBoard(c)
	}
}

func (c *Permanent) CanBlock(attacker *Permanent) bool {
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

func (c *Permanent) HandleComingIntoPlay() {
	if c.Owner == nil {
		panic("permanent has unset owner")
	}
	if c.Bloodthirst > 0 && c.Owner.Opponent().DamageThisTurn > 0 {
		c.Plus1Plus1Counters += c.Bloodthirst
	}
	if c.EntersPlayEffect != nil {
		c.Owner.ResolveEffect(c.EntersPlayEffect)
	}
}

/*
	Most creatures don't do anything special when they deal damage.
	Currently just ones with Lifelink do something extra.
*/
func (c *Permanent) DidDealDamage(damage int) {
	if c.Lifelink && damage > 0 {
		c.Owner.Life += damage
	}
}

func (c *Permanent) ActivateAbility(cost *Cost, target *Permanent) {
	if c.ActivatedAbility == nil {
		panic("tried to activate a permanent without an ability")
	}
	c.ActivatedThisTurn = true
	selectedForCost := cost.Effect.SelectedForCost
	if c.ActivatedAbility.Cost.Effect.EffectType == ReturnToHand {
		selectedForCost.Owner.RemoveFromBoard(selectedForCost)
		selectedForCost.Owner.Hand = append(selectedForCost.Owner.Hand, selectedForCost.Card.Name)
	}
	if c.ActivatedAbility.EffectType == Untap {
		target.Tapped = false
	}
}
