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


func (c *Card) AsciiImage() {
	cardWidth := 11
	cardHeight := 5
	imageGrid := [cardHeight][cardWidth]int
	for y := 0; y < cardHeight; y++ {
		for x := 0; x < cardWidth; x++ {
				if x == 0 or x == cardWidth - 1:
					imageGrid[-1].append('|')
				elif y == 0 or y == cardHeight - 1:
					imageGrid[-1].append('-')
				else:
					imageGrid[-1].append(' ')
		}
	}

	initialIndex := 2

	if showBack:
		middleX := cardWidth / 2
		middleY := cardHeight / 2

		noon := (middleX, middleY - 1, '*')
		two := (middleX + 2, middleY, '*')
		ten := (middleX - 2, middleY, '*')
		seven := (middleX - 1, middleY + 1, '*')
		four := (middleX + 1, middleY + 1, '*')

		points := [noon, two, four, seven, ten]
		for _, p := range points {
			imageGrid[p[1]][p[0]] = p[2]
		}
	
	if !showBack {
		ccRow = 1
		ccString = fmt.Sprintf("%v", c.ManaCost)
		for x in range(initialIndex, len(ccString) + initialIndex):
			imageGrid[ccRow][x] = ccString[x-initialIndex]
	}	

	if !showBack:
		nameRow = 2
		words := fmt.Sprintf("%v", c.Name)
		for _, word := range words {
			wordWidth = min(3, len(word))
			if len(words) == 1:
				wordWidth = min(len(word), cardWidth - 4)
			for x in range(initialIndex, wordWidth + initialIndex):
				imageGrid[nameRow][x] = word[x-initialIndex]
			initialIndex += wordWidth + 1
			if initialIndex >= cols - wordWidth - 1:
				break

	if not showBack:
		if c.IsCreature {
			initialIndex = 2
			statsRow = 3
			statsString = fmt.Sprintf("%v/%v", c.Power, c.Toughnesss)
			for x in range(initialIndex, len(statsString) + initialIndex):
				image_grid[statsRow][x] = statsString[x-initialIndex]
		}

	if not show_back:
		if c.IsTapped {
		if Card.tapped(card_state):
			tapped_row = 0
			initial_index = 0
			tapped_string = "TAPPED"
			for x in range(initial_index, len(tapped_string) + initial_index):
				image_grid[tapped_row][x] = tapped_string[x-initial_index]

	return image_grid
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