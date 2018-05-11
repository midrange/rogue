package game

import (
	"fmt"
	"log"
	"math/rand"
    "strings"
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
	for _, arrayLine := range c.AsciiImage(false) {
		for _, char := range arrayLine {
			fmt.Printf(char)
		}
		fmt.Printf("%v", "\n")
	}
}


func (c *Card) AsciiImage(showBack bool) [CARD_HEIGHT][CARD_WIDTH]string {
	const cardWidth = CARD_WIDTH
	const cardHeight = CARD_HEIGHT
	imageGrid := [cardHeight][cardWidth]string{}
	for y := 0; y < cardHeight; y++ {
		for x := 0; x < cardWidth; x++ {
			if x == 0 || x == cardWidth - 1 {
				imageGrid[y][x] = string('|')
			} else if y == 0 || y == cardHeight - 1 {
				imageGrid[y][x] = string('-')
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
				wordWidth = Min(len(word), cardWidth - 4)
			}
			for x := initialIndex; x < wordWidth + initialIndex; x++ {
				imageGrid[nameRow][x] = string(word[x-initialIndex])
			}
			initialIndex += wordWidth + 1
			if initialIndex >= cardWidth - wordWidth - 1 {
				break
			}
		}

		if c.IsCreature {
			initialIndex := 2
			statsRow := 3
			statsString := fmt.Sprintf("%v/%v", c.Power, c.Toughness)
			for x := initialIndex; x < len(statsString) + initialIndex; x++ {
				imageGrid[statsRow][x] = string(statsString[x-initialIndex])
			}

			ccRow := 1
			ccString := fmt.Sprintf("%v", c.ManaCost)
			for x := initialIndex; x < len(ccString) + initialIndex; x++ {
				imageGrid[ccRow][x] = string(ccString[x-initialIndex])
			}
		}

		if c.Tapped {
			tappedRow := 0
			initialIndex := 0
			tappedString := "TAPPED"
			for x := initialIndex; x < len(tappedString) + initialIndex; x++ {
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