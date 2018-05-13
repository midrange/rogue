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
		g.passTurn()
		g.passTurn()
	}

	if g.IsOver() {
		t.Fatalf("game was unexpectedly over")
	}

	g.passTurn()
	if !g.Players[1].Lost() {
		t.Fatalf("the player on the draw should have lost by decking")
	}
}

func TestTwoBearsFighting(t *testing.T) {
	g := NewGame(topBear(), topBear())

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

}
