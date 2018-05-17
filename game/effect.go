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

	// from Card
	CastingCost        *CastingCost // when an Effect is a kicker, it has a Cost
	Hexproof           bool
	Kicker             *Effect
	Plus1Plus1Counters int
	Power              int
	Summon             *Card
	Toughness          int
	Untargetable       bool

	// from Action
	ActionType    ActionType
	Target        *Permanent
	With          *Permanent
	WithKicker    bool
	WithPhyrexian bool
}

func NewEffect(action *Action) *Effect {
	card := action.Card
	effect := card.Effect
	effect.Kicker = card.Kicker

	effect.ActionType = action.Type
	effect.Target = action.Target
	effect.With = action.With
	effect.WithKicker = action.WithKicker
	effect.WithPhyrexian = action.WithPhyrexian
	return effect
}
