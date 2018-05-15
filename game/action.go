package game

import (
	"fmt"
)

type Action struct {
	Type          ActionType
	Card          *Card
	Target        *Card
	WithKicker    bool
	WithPhyrexian bool
}

type ActionType int

const (
	Pass ActionType = iota
	Play
	DeclareAttack
	Attack
	Block
	UseForMana
	ChooseTargetAndMana
)

func (a *Action) pronoun() string {
	if a.Target.Owner == a.Card.Owner {
		return "your"
	}
	return "their"
}

// For debugging and logging. Don't use this in the critical path.
func (a *Action) String() string {
	switch a.Type {
	case Pass:
		return "pass"
	case ChooseTargetAndMana:
		fallthrough
	case Play:
		if a.WithKicker {
			if a.Target == nil {
				return fmt.Sprintf("%v: %v with kicker", a.Card.Kicker.CastingCost.Colorless, a.Card.String())
			}
			return fmt.Sprintf("%v: %v on %v %v with kicker", a.Card.Kicker.CastingCost.Colorless, a.Card.String(), a.pronoun(), a.Target.String())
		}
		if a.Card.IsLand {
			return fmt.Sprintf("%v", a.Card.String())
		}
		if a.Target == nil {
			return fmt.Sprintf("%v: %v", a.Card.CastingCost.Colorless, a.Card.String())
		}
		return fmt.Sprintf("%v: %v on %v %v", a.Card.CastingCost.Colorless, a.Card.String(), a.pronoun(), a.Target.String())
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
