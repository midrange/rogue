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
	g := NewGame(
		deckWithTopAndForests(GrizzlyBears),
		deckWithTopAndForests(GrizzlyBears))

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

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
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
		With:   defendingBear,
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
	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
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

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(DeclareBlockers)
	if len(g.Actions(false)) > 1 {
		t.Fatal("expected no legal blocks")
	}
}

func TestSilhanaCantBeTargeted(t *testing.T) {
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

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)

	g.playCreature()

	if g.Priority().GetBoard()[2].Power() != 2 {
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

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)

	g.playCreature()
	g.passTurn()

	g.passTurn()

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	if len(g.Actions(false)) > 2 {
		t.Fatal("expected the small skulk couldnt block the big skulk")
	}
}

func TestVaultSkirgeLoseAndGain(t *testing.T) {
	g := NewGame(deckWithTopAndForests(VaultSkirge), deckWithTopAndForests(VaultSkirge))
	g.playLand()
	g.playCreaturePhyrexian()
	if g.Priority().Life != 18 {
		g.Print()
		t.Fatal("expected the player to lose 2 life from Vault Skirge casting")
	}
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(Main2)
	if g.Priority().Life != 19 {
		t.Fatal("expected the player to gain 1 life from Vault Skirge attack")
	}
}

func TestNestInvader(t *testing.T) {
	g := NewGame(deckWithTopAndForests(NestInvader), deckWithTopAndForests(NestInvader))
	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playCreature()

	if len(g.Priority().Creatures()) != 2 {
		t.Fatal("expected the player to have a Nest Invader and a token")
	}

	g.playManaAbilityAction()

	if g.Priority().ColorlessManaPool != 1 {
		t.Fatal("expected the player to have a colorless floating")
	}

	if len(g.Priority().Creatures()) != 1 {
		t.Fatal("expected the player to have a Nest Invader, with the token now dead")
	}
}

func TestBurningTreeEmissary(t *testing.T) {
	g := NewGame(deckWithTopAndForests(BurningTreeEmissary), deckWithTopAndForests(BurningTreeEmissary))
	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playCreature()

	if g.Priority().ColorlessManaPool != 2 {
		t.Fatal("expected the player to have 2 mana from BurningTreeEmissary")
	}
}

func TestQuirionRanger(t *testing.T) {
	g := NewGame(deckWithTopAndForests(QuirionRanger), deckWithTopAndForests(QuirionRanger))
	g.playLand()
	g.playCreature()
	g.playActivatedAbility()

	if len(g.Priority().Board) != 1 {
		t.Fatal("expected the player to only have Quirion Ranger in play, not ", g.Priority().Board)
	}

	if len(g.Priority().Hand) != 6 {
		t.Fatal("expected player to have 6 cards in hand after returning Forest with Quirion Ranger")
	}
}

func TestElephantGuide(t *testing.T) {
	skirgeGuide := NewEmptyDeck()
	skirgeGuide.Add(1, ElephantGuide)
	skirgeGuide.Add(1, VaultSkirge)
	skirgeGuide.Add(58, Forest)

	skirgeGuide2 := NewEmptyDeck()
	skirgeGuide2.Add(1, ElephantGuide)
	skirgeGuide2.Add(1, VaultSkirge)
	skirgeGuide2.Add(58, Forest)

	g := NewGame(skirgeGuide, skirgeGuide2)

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playAura()
	g.passTurn()

	g.playLand()
	g.playAura()
	g.passTurn()

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(DeclareBlockers)
	g.doBlockAction()
	g.passUntilPhase(Main2)

	if len(g.Priority().Board) != 4 {
		t.Fatal("expected the attacker to have gotten a token")
	}

	if len(g.Priority().Opponent().Board) != 4 {
		t.Fatal("expected the opponent to have gotten a token")
	}
}

func TestRancor(t *testing.T) {
	skirgeRancor := NewEmptyDeck()
	skirgeRancor.Add(1, Rancor)
	skirgeRancor.Add(1, VaultSkirge)
	skirgeRancor.Add(58, Forest)

	skirgeRancor2 := NewEmptyDeck()
	skirgeRancor2.Add(1, Rancor)
	skirgeRancor2.Add(1, VaultSkirge)
	skirgeRancor2.Add(58, Forest)

	g := NewGame(skirgeRancor, skirgeRancor2)

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playAura()

	if len(g.Priority().Hand) != 5 {
		t.Fatal("expected the hand size to be 5 before rancor return")
	}

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(DeclareBlockers)
	g.doBlockAction()

	g.passUntilPhase(Main2)

	if len(g.Priority().Hand) != 6 {
		t.Fatal("expected the rancor to return to hand")
	}
}

func TestFaerieMiscreant(t *testing.T) {
	twoMiscreants := NewEmptyDeck()
	twoMiscreants.Add(2, FaerieMiscreant)
	twoMiscreants.Add(58, Island)
	g := NewGame(twoMiscreants, deckWithTopAndForests(BurningTreeEmissary))

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.passTurn()

	g.playCreature()

	if len(g.Priority().Hand) != 6 {
		t.Fatal("expected the player to have 6 cards after playing 2nd miscreant")
	}
}

func TestMutagenicGrowth(t *testing.T) {
	skirgeGrowth := NewEmptyDeck()
	skirgeGrowth.Add(1, MutagenicGrowth)
	skirgeGrowth.Add(1, VaultSkirge)
	skirgeGrowth.Add(58, Forest)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(skirgeGrowth, allForests)

	g.playLand()
	g.playCreature()
	g.playInstant()

	if len(g.Priority().Hand) != 4 {
		t.Fatal("expected the hand size to be 4 after forest, skirge, mutagenic growth")
	}
}

func TestGush(t *testing.T) {
	gush := NewEmptyDeck()
	gush.Add(1, Gush)
	gush.Add(59, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(gush, allForests)

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()
	g.playInstant()

	if len(g.Priority().Hand) != 9 {
		t.Fatal("expected the hand size to be 9 after Gush: draw 7, island, draw, island, play gush, draw 2")
	}
}

func TestSnap(t *testing.T) {
	snapSkirge := NewEmptyDeck()
	snapSkirge.Add(1, Snap)
	snapSkirge.Add(1, VaultSkirge)
	snapSkirge.Add(59, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(snapSkirge, allForests)

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.playInstant()

	if len(g.Priority().Hand) != 5 {
		t.Fatal("expected the hand size to be 5 after Snap: draw 7, island, skirge, draw, island, snap")
	}

	for _, land := range g.Priority().Lands() {
		if land.Tapped {
			t.Fatal("expected all player's lands to be untapped from snap")
		}
	}
}

func TestCounterspell(t *testing.T) {
	skirge := NewEmptyDeck()
	skirge.Add(1, VaultSkirge)
	skirge.Add(59, Island)

	counter := NewEmptyDeck()
	counter.Add(1, Counterspell)
	counter.Add(59, Island)

	g := NewGame(counter, skirge)

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.putCreatureOnStackAndPass()

	if len(g.Stack) != 1 {
		t.Fatal("expected there to be Vault Skirge on the stack")
	}
	g.playInstant()
	if len(g.Stack) != 0 {
		g.Print()
		t.Fatal("expected there to be no spells on the stack after Counterspell")
	}
}

func TestDazeNotPaid(t *testing.T) {
	skirge := NewEmptyDeck()
	skirge.Add(1, VaultSkirge)
	skirge.Add(59, Island)

	daze := NewEmptyDeck()
	daze.Add(1, Daze)
	daze.Add(59, Island)

	g := NewGame(daze, skirge)

	g.playLand()
	g.passTurn()

	g.playLand()
	g.putCreatureOnStackAndPass()

	if len(g.Stack) != 1 {
		t.Fatal("expected there to be Vault Skirge on the stack")
	}
	g.playInstant()
	g.TakeAction(&Action{
		Type:        MakeChoice,
		AfterEffect: &Effect{EffectType: Countermagic, SpellTargetId: g.ChoiceEffect.SpellTargetId},
	})
	if len(g.Stack) != 0 {
		t.Fatal("expected there to be no spells on the stack after Daze")
	}
}

func TestDazePaid(t *testing.T) {
	skulk := NewEmptyDeck()
	skulk.Add(1, SkarrganPitskulk)
	skulk.Add(59, Island)

	daze := NewEmptyDeck()
	daze.Add(1, Daze)
	daze.Add(59, Island)

	g := NewGame(skulk, daze)

	g.playLand()
	g.passTurn()

	g.playLand()
	g.passTurn()

	g.playLand()
	g.putCreatureOnStackAndPass()

	if len(g.Stack) != 1 {
		t.Fatal("expected there to be Vault Skirge on the stack")
	}
	g.playInstant()

	// pay for Daze
	for _, a := range g.Actions(false) {
		if a.AfterEffect.EffectType == TapLand {
			g.TakeAction(a)
			break
		}
	}

	g.TakeAction(&Action{Type: PassPriority})
	if len(g.Creatures()) != 1 {
		t.Fatal("expected there to be Vault Skirge in play after Daze was paid")
	}
	if len(g.Lands()) != 2 {
		t.Fatal("expected there to be 2 in play after Daze")
	}
}

func TestSpellstutterSpriteSucceeds(t *testing.T) {
	sprite := NewEmptyDeck()
	sprite.Add(1, SpellstutterSprite)
	sprite.Add(1, NettleSentinel)
	sprite.Add(58, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(sprite, allForests)

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()

	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.Name == NettleSentinel {
			g.TakeAction(a)
			break
		}
	}

	g.putCreatureOnStackAndPass()

	if len(g.Stack) != 2 {
		t.Fatal("expected two creatures on the stack")
	}

	g.TakeAction(&Action{Type: PassPriority})

	if len(g.Stack) != 2 {
		t.Fatal("expected creature and spellstutter on stack")
	}

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: PassPriority})

	if len(g.Stack) != 0 {
		t.Fatal("expected 0 creatures on the stack ", g.Stack)
	}

	if len(g.Priority().Board) != 4 {
		t.Fatal("expected no NettleSentinel in play ", g.Priority().Board)
	}
}

func TestSpellstutterSpriteFails(t *testing.T) {
	sprite := NewEmptyDeck()
	sprite.Add(1, VaultSkirge)
	sprite.Add(1, SpellstutterSprite)
	sprite.Add(58, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(sprite, allForests)

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()
	g.passTurn()

	g.passTurn()

	g.playLand()

	for _, a := range g.Priority().PlayActions(true, false) {
		if a.Card != nil && a.Card.Name == VaultSkirge {
			g.TakeAction(a)
			break
		}
	}
	g.putCreatureOnStackAndPass()

	if len(g.Stack) != 2 {
		t.Fatal("expected two creatures on the stack")
	}

	g.TakeAction(&Action{Type: PassPriority})

	if len(g.Stack) != 2 {
		t.Fatal("expected creature and spellstutter on stack")
	}

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: PassPriority})

	if len(g.Stack) != 1 {
		t.Fatal("expected vault skirge still to be on the stack, not enough faeries ", g.Stack)
	}
}

func TestPonder(t *testing.T) {
	ponder := NewEmptyDeck()
	ponder.Add(1, Ponder)
	ponder.Add(59, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(ponder, allForests)
	g.playLand()

	g.playSorcery()

	// DecideOnPonder action
	g.TakeAction(g.Actions(false)[0])

	if len(g.Priority().Hand) != 6 {
		panic("expected 6 cards in hand after Ponder")
	}
}

func TestPreordain(t *testing.T) {
	preordain := NewEmptyDeck()
	preordain.Add(1, Preordain)
	preordain.Add(59, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(preordain, allForests)

	g.playLand()

	g.playSorcery()

	if len(g.Priority().Hand) != 5 {
		g.Print()
		panic("expected 5 cards in hand on cast Preordain")
	}
	// Scry action
	g.TakeAction(g.Actions(false)[0])

	if len(g.Priority().Hand) != 6 {
		panic("expected 6 cards in hand after Preordain resolved")
	}
}

func TestNinjaOfTheDeepWater(t *testing.T) {
	ninja := NewEmptyDeck()
	ninja.Add(1, NinjaOfTheDeepHours)
	ninja.Add(1, FaerieMiscreant)
	ninja.Add(59, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(ninja, allForests)

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.passTurn()

	g.playLand()

	g.TakeAction(&Action{Type: PassPriority})
	g.TakeAction(&Action{Type: Pass})
	g.attackWithEveryone()
	g.passUntilPhase(CombatDamage)

	g.playCreature()
	g.passUntilPhase(Main2)

	if g.Defender().Life != 18 {
		panic("expected defender's life to be 18 after Ninja attack")
	}

	if len(g.Attacker().Hand) != 6 {
		g.Print()
		panic("expected defender's hand to have 6 cards after Ninja attack")
	}
}

func TestDelverFlips(t *testing.T) {
	delver := NewEmptyDeck()
	delver.Add(1, DelverOfSecrets)
	delver.Add(6, Island)
	delver.Add(1, Ponder)
	delver.Add(52, Island)

	allForests := NewEmptyDeck()
	allForests.Add(60, Forest)

	g := NewGame(delver, allForests)

	g.playLand()
	g.playCreature()
	g.passTurn()

	g.passUntilPhase(Draw)
	g.TakeAction(g.Actions(false)[0])

	if g.Attacker().Creatures()[0].Name != InsectileAberration {
		panic("expected Delver to transform")
	}

}
