/*
	An Effect can be created by a spell, a triggered ability, or an
	activated ability.

	An Effect can be the property of a Card. For example, when the
	Card is a Giant Growth, its effect property would be
	&Effect{power:3, toughness:3}.

	An Effect can also be the Morbid or Kicker property of a Card, to designate a
	Effect that only happens under special circumstances.

*/

package game

import ()

type Effect struct {

	// when an Effect is a kicker, it has a Cost
	CastingCost *CastingCost

	// these properties modify a Card the Effect targets
	Hexproof           bool
	Kicker             *Effect
	Plus1Plus1Counters int
	Power              int
	Summon             *Card
	Toughness          int
	Untargetable       bool

	// these properties get copied from the Action object from which the Effect is created
	Source *Permanent
	Target *Permanent
}

func NewEffect(action *Action) *Effect {
	card := action.Card
	effect := card.Effect
	if action.WithKicker {
		effect.Kicker = card.Kicker
	}
	effect.Source = action.Source
	effect.Target = action.Target
	return effect
}
