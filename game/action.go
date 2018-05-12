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
	Pass ActionType = iota
	Play
	Attack
	Block
)

// For debugging and logging. Don't use this in the critical path.
func (a *Action) String() string {
	switch a.Type {
	case Pass:
		return "pass"
	case Play:
		return fmt.Sprintf("play %v", a.Card.String())
	case Attack:
		return fmt.Sprintf("attack with %v", a.Card.String())
	case Block:
		return fmt.Sprintf("%v blocks %v", a.Card.String(), a.Target.String())
	}
	panic("control should not reach here")
}
