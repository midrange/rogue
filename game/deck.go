package game

import (
	"math/rand"
)

type Deck struct {
	Cards        []CardName
	FailedToDraw bool
}

func NewEmptyDeck() *Deck {
	return &Deck{
		Cards: []CardName{},
	}
}

// Constructs a new deck from the count of each card.
func NewDeck(decklist map[CardName]int) *Deck {
	deck := NewEmptyDeck()
	for name, count := range decklist {
		deck.Add(count, name)
	}
	deck.Shuffle()
	return deck
}

func Stompy() *Deck {
	return NewDeck(map[CardName]int{
		BurningTreeEmissary: 4,
		ElephantGuide:       4,
		Forest:              17,
		HungerOfTheHowlpack: 4,
		NestInvader:         4,
		NettleSentinel:      4,
		// QuirionRanger:       4,
		Rancor:             4,
		SkarrganPitskulk:   4,
		SilhanaLedgewalker: 3,
		VaultSkirge:        4,
		VinesOfVastwood:    4,
	})
}

func MonoBlueDelver() *Deck {
	return NewDeck(map[CardName]int{
		// DelverOfSecrets:     4,
		FaerieMiscreant:     4,
		SpellstutterSprite:  4,
		Island:              18,
		NinjaOfTheDeepHours: 4,
		MutagenicGrowth:     4,
		Ponder:              4,
		Preordain:           4,
		Counterspell:        4,
		Daze:                4,
		Snap:                4,
		Gush:                2,
	})
}

func (d *Deck) Draw() CardName {
	if len(d.Cards) == 0 {
		d.FailedToDraw = true
		return NoCard
	}
	answer := d.Cards[0]
	d.Cards = d.Cards[1:]
	return answer
}

func (d *Deck) Shuffle() {
	for i := len(d.Cards) - 1; i > 0; i-- {
		// Swap the ith card with a random one in [0..i]
		j := rand.Int() % (i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Adds cards to the deck, on bottom
func (d *Deck) Add(n int, name CardName) {
	for i := 0; i < n; i++ {
		d.Cards = append(d.Cards, name)
	}
}

// Adds cards to the deck, on top
func (d *Deck) AddToTop(n int, name CardName) {
	for i := 0; i < n; i++ {
		d.Cards = append([]CardName{name}, d.Cards...)
	}
}
