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
	DeclareAttack
	Attack
	Block
)

// For debugging and logging. Don't use this in the critical path.
func (a *Action) String() string {
	switch a.Type {
	case PassPriority:
		return "pass priority"
	case PassTurn:
		return "pass turn"
	case Play:
		return fmt.Sprintf("play %v", a.Card.String())
	case DeclareAttack:
		return "enter attack step"
	case Attack:
		return fmt.Sprintf("attack with %v", a.Card.String())
	case Block:
		return fmt.Sprintf("%v blocks %v", a.Card.String(), a.Target.String())
	}
	panic("control should not reach here")
}
