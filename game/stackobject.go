package game

import (
	"fmt"
)

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

func (s *StackObject) String() string {
	if s.Card != nil {
		return fmt.Sprintf("resolve %s", s.Card)
	}
	return fmt.Sprintf("resolve ability from %s", s.Source)
}
