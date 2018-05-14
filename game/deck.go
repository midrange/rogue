package game

import (
	"math/rand"
)

type Deck struct {
	Cards        []*Card
	FailedToDraw bool
}

func NewEmptyDeck() *Deck {
	return &Deck{
		Cards: []*Card{},
	}
}

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
		Forest:              20,
		HungerOfTheHowlpack: 5,
		NettleSentinel:      15,
		Rancor:              5,
		SkarrganPitskulk:    5,
		SilhanaLedgewalker:  10,
		VinesOfVastwood:     5,
	})
}

func (d *Deck) Draw() *Card {
	if len(d.Cards) == 0 {
		d.FailedToDraw = true
		return nil
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
		d.Cards = append(d.Cards, NewCard(name))
	}
}
