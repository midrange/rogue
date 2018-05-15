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

func TestDecking(t *testing.T) {
	g := NewGame(deckWithTopAndForests(GrizzlyBears), deckWithTopAndForests(GrizzlyBears))

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
	g := NewGame(deckWithTopAndForests(GrizzlyBears), deckWithTopAndForests(GrizzlyBears))

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
	g := NewGame(topNettleVines(), deckWithTopAndForests(GrizzlyBears))

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
	g := NewGame(topNettleVines(), deckWithTopAndForests(GrizzlyBears))

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

// A deck stacked with a NettleSentinel and two VinesOfVastwood on top and all the rest forests
func topNettleVines() *Deck {
	deck := NewEmptyDeck()
	deck.Add(2, VinesOfVastwood)
	deck.Add(1, NettleSentinel)
	deck.Add(57, Forest)
	return deck
}

func TestSilhanasDontMeet(t *testing.T) {
	g := NewGame(deckWithTopAndForests(SilhanaLedgewalker), deckWithTopAndForests(SilhanaLedgewalker))

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
	g := NewGame(deckWithTopAndForests(SilhanaLedgewalker), deckWithTopAndForests(VinesOfVastwood))

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

func TestSkarrganPitskulkBloodthirst(t *testing.T) {
	twoSkulksDeck := NewEmptyDeck()
	twoSkulksDeck.Add(2, SkarrganPitskulk)
	twoSkulksDeck.Add(58, Forest)
	g := NewGame(twoSkulksDeck, deckWithTopAndForests(SkarrganPitskulk))

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)

	g.playCreature()

	if g.Priority().Board[2].Power() != 2 {
		t.Fatal("expected a bloodthirsted skulk")
	}

}

func TestSkarrganPitskulksDontMeet(t *testing.T) {
	twoSkulksDeck := NewEmptyDeck()
	twoSkulksDeck.Add(2, SkarrganPitskulk)
	twoSkulksDeck.Add(58, Forest)
	g := NewGame(twoSkulksDeck, deckWithTopAndForests(SkarrganPitskulk))

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)

	g.playCreature()
	g.passTurn()

	g.passTurn()

	g.TakeAction(&Action{Type: DeclareAttack})
	g.attackWithEveryone()
	if len(g.Actions(false)) > 2 {
		t.Fatal("expected the small skulk couldnt block the big skulk")
	}

}
