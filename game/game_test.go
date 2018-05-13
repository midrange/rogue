package game

import (
	"testing"
)

// A deck stacked with a bear on top and all the rest forests
func topBear() *Deck {
	deck := NewEmptyDeck()
	deck.Add(1, GrizzlyBears)
	deck.Add(59, Forest)
	return deck
}

func TestDecking(t *testing.T) {
	g := NewGame(topBear(), topBear())

	// When each player passes the turn 53 times, both decks should be out of cards
	for i := 0; i < 53; i++ {
		g.PassTurn()
		g.PassTurn()
	}

	if g.IsOver() {
		t.Fatalf("game was unexpectedly over")
	}

	g.PassTurn()
	if !g.Players[1].Lost() {
		t.Fatalf("the player on the draw should have lost by decking")
	}
}

func TestFighting(t *testing.T) {
	NewGame(topBear(), topBear())

	// TODO: have two bears fight
}
