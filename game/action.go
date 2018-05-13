package game

import (
	"fmt"
)

type Action struct {
	Type   ActionType
	Card   *Card
	Target *Card
}

type ActionType int

const (
	PassPriority ActionType = iota
	PassTurn
	Play
	PlayWithKicker
	DeclareAttack
	Attack
	Block
	UseForMana
)

// For debugging and logging. Don't use this in the critical path.
func (a *Action) String() string {
	switch a.Type {
	case PassPriority:
		return "pass priority"
	case PassTurn:
		return "pass turn"
	case Play:
		if a.Target == nil {
			return fmt.Sprintf("play %v", a.Card.String())
		}
		return fmt.Sprintf("play %v on %v", a.Card.String(), a.Target.String())
	case PlayWithKicker:
		if a.Target == nil {
			return fmt.Sprintf("play %v with kicker", a.Card.String())
		}
		return fmt.Sprintf("play %v on %v with kicker", a.Card.String(), a.Target.String())
	case DeclareAttack:
		return "enter attack step"
	case Attack:
		return fmt.Sprintf("attack with %v", a.Card.String())
	case Block:
		return fmt.Sprintf("%v blocks %v", a.Card.String(), a.Target.String())
	case UseForMana:
		return fmt.Sprintf("tap %v for mana", a.Card.String())
	}
	panic("control should not reach here")
}
