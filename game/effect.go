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
	// when an Effect is a kicker, it has a CastingCost
	CastingCost *CastingCost

	// these properties modify a Permanent the Effect targets
	Colorless          int
	Hexproof           bool
	Plus1Plus1Counters int
	Power              int
	Toughness          int
	Untargetable       bool

	/*
		gets sets when the action for the Effect is withKicker
		this child Effect adds on top of the normal Effect
	*/
	Kicker *Effect

	// sometimes an effect summons a creature
	Summon CardName

	// Source is the source of activated abilities, nil for other effects.
	Source *Permanent

	Target *Permanent

	// for effects from targeted spells
	EffectType EffectType
	TargetType *TargetType

	// when an effect has a cost, e.g. a Quirion Ranger
	Cost *Effect
}

//go:generate stringer -type=EffectType
type EffectType int

const (
	ReturnToHand EffectType = iota
	Untap
)

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
