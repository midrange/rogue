package game

import (
	"fmt"
)

type StackObjectId int

const NoStackObjectId StackObjectId = 0

type StackObject struct {
	Type                            ActionType
	Card                            *Card // for spell-based stack objects
	Cost                            *Cost
	EntersTheBattleFieldSpellTarget StackObjectId
	Id                              StackObjectId
	Kicker                          *Effect
	Player                          PlayerId
	Selected                        []PermanentId
	Source                          *Permanent
	SpellTarget                     StackObjectId
	Target                          *Permanent // a target that is a Permanent (players not yet handled)
	WithNinjitsu                    bool
}

func (s *StackObject) String() string {
	if s.Type == EntersTheBattlefieldEffect {
		return fmt.Sprintf("%s Enters the Battlefield effect", s.Card)
	}
	if s.Card != nil {
		return fmt.Sprintf("resolve %s", s.Card)
	}
	return fmt.Sprintf("resolve ability from %s", s.Source)
}
