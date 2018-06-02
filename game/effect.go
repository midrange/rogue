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
	// required for effect to occur
	Condition *Condition

	// when an Effect is a kicker, it has a Cost
	Cost *Cost

	// these properties modify a Permanent the Effect targets, or the Game state
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

	SelectedForCost *Permanent
	SpellTarget     *StackObject
	Target          *Permanent

	// for effects from targeted spells
	EffectType EffectType
	Selector   *Selector

	// for non-targeted effects of spells, such as Snap
	Selected []*Permanent
}

//go:generate stringer -type=EffectType
type EffectType int

const (
	AddMana EffectType = iota
	Countermagic
	DrawCard
	ManaSink
	ReturnToHand
	Untap
)

func UpdatedEffectForStackObject(stackObject *StackObject, effect *Effect) *Effect {
	newEffect := effect
	newEffect.Kicker = stackObject.Kicker
	newEffect.Source = stackObject.Source
	newEffect.Target = stackObject.Target
	newEffect.Selected = stackObject.Selected
	newEffect.SpellTarget = stackObject.SpellTarget
	return newEffect
}
