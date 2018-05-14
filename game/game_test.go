package game

import (
	"testing"
)

// A deck stacked with a certain card c on top and all the rest forests
func deckWithTopAndForests(name CardName) *Deck {
	deck := NewEmptyDeck()
	deck.Add(1, name)
	deck.Add(59, Forest)
	return deck
}

func topBear() *Deck {
	return deckWithTopAndForests(GrizzlyBears)
}

func topSilhana() *Deck {
	return deckWithTopAndForests(SilhanaLedgewalker)
}

func topVines() *Deck {
	return deckWithTopAndForests(VinesOfVastwood)
}

// A deck stacked with a NettleSentinel and two VinesOfVastwood on top and all the rest forests
func topNettleVines() *Deck {
	deck := NewEmptyDeck()
	deck.Add(2, VinesOfVastwood)
	deck.Add(1, NettleSentinel)
	deck.Add(57, Forest)
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

	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	attackingBear := g.Attacker().GetCreature(GrizzlyBears)
	if !attackingBear.Attacking {
		t.Fatal("expected attacking bear to be attacking")
	}
	defendingBear := g.Defender().GetCreature(GrizzlyBears)
	if defendingBear.Attacking {
		t.Fatal("expected defending bear to not be attacking")
	}
	g.TakeAction(&Action{
		Type:   Block,
		Card:   defendingBear,
		Target: attackingBear,
	})
	g.passUntilPhase(Main2)

	if len(g.Attacker().Creatures()) != 0 {
		t.Fatal("expected attacking bear to die")
	}

	if len(g.Defender().Creatures()) != 0 {
		t.Fatal("expected defending bear to die")
	}
}

func TestVinesOfVastwoodBuff(t *testing.T) {
	g := NewGame(topNettleVines(), topBear())

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playKickedInstant()
	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)

	if g.Defender().Life != 14 {
		t.Fatal("expected defender life total to be 14, instead was ", g.Defender().Life)
	}
}

func TestVinesOfVastwoodUntargetable(t *testing.T) {
	g := NewGame(topNettleVines(), topBear())

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playInstant()

	actions := g.Actions(false)

	for _, action := range actions {
		if action.Type == Play {
			t.Fatal("expected no legal Plays for the 2nd Vines of Vastwood")
		}

	}
}

func TestSilhanasDontMeet(t *testing.T) {
	g := NewGame(topSilhana(), topSilhana())

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

	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	g.passUntilPhase(DeclareBlockers)
	if len(g.Actions(false)) > 1 {
		t.Fatal("expected no legal blocks")
	}
}

func TestSilhanaCantBeTargetted(t *testing.T) {
	g := NewGame(topSilhana(), topVines())

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	for _, action := range g.Actions(false) {
		if action.Type == Play {
			t.Fatal("expected no legal targets")
		}

	}

}
