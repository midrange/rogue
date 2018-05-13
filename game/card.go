package game

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
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
	Auras      []*Card
	Tapped     bool
	TurnPlayed int

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
	BaseTrample   bool
}

//go:generate stringer -type=CardName
type CardName int

const (
	Forest CardName = iota
	GrizzlyBears
	Rancor
	NettleSentinel
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
		return &Card{
			IsCreature:    true,
			BasePower:     2,
			BaseToughness: 2,
			ManaCost:      1,
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

func (c *Card) String() string {
	return fmt.Sprintf("%v (%v) %v/%v", c.Name, c.ManaCost, c.Power(), c.Toughness())
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

			ccRow := 1
			ccString := fmt.Sprintf("%v", c.ManaCost)
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
	if c.Name == Forest && !c.Tapped {
		[]*Action{&Action{Type: UseForMana}}
	}
	return []*Action{}
}
