package game

import (
	"fmt"
	"log"
	"math/rand"
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
	cardWidth := 10
	printCardBorder(cardWidth)
	printBlankCardLine(cardWidth)
	if c.IsCreature {
		printCardTextLine(cardWidth, fmt.Sprintf("%v", c.ManaCost))
	} else {
		printBlankCardLine(cardWidth)
	}
	printCardTextLine(cardWidth, fmt.Sprintf("%v", c.Name))
	if c.IsCreature {
		printCardTextLine(cardWidth, fmt.Sprintf("%v/%v", c.Power, c.Toughness))
	} else {
		printBlankCardLine(cardWidth)
	}
	printBlankCardLine(cardWidth)
	printCardBorder(cardWidth)
}

func printCardBorder(cardWidth int) {
	for i := 0; i < cardWidth; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
}

func printCardTextLine(cardWidth int, name string) {
	nameLine := "| "
	nameIndex := 0
	for {
		if nameIndex < len(name) {
			nameLine += string(name[nameIndex])
		} else {
			nameLine += " "
		}
		nameIndex += 1
		if len(nameLine) >= cardWidth - 2 { break; }
	}
	nameLine += " |"
	fmt.Printf("%v\n", nameLine)
}

func printBlankCardLine(cardWidth int) {
	fmt.Printf("|",)
	for i := 0; i < cardWidth - 2; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("|\n")
}