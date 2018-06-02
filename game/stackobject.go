package game

import ()

type StackObject struct {
	Type                            ActionType
	Card                            *Card // for spell-based stack objects
	Cost                            *Cost
	EntersTheBattleFieldSpellTarget *StackObject
	Kicker                          *Effect
	Player                          *Player
	Selected                        []*Permanent
	Source                          *Permanent
	SpellTarget                     *StackObject
	Target                          *Permanent // a target that is a Permanent (players not yet handled)
	WithNinjitsu                    bool
}
