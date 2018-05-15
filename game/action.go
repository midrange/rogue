package game

import (
	"fmt"
)

type Action struct {
	Type       ActionType
	Card       *Card
	With       *Permanent
	Target     *Permanent
	WithKicker bool
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

func (a *Action) targetPronoun(p *Player) string {
	if a.Target.Owner == p {
		return "your"
	}
	return "their"
}

// For debugging and logging. Don't use this in the critical path.
func (a *Action) ShowTo(p *Player) string {
	switch a.Type {
	case Pass:
		return "pass"
	case ChooseTargetAndMana:
		fallthrough
	case Play:
		if a.WithKicker {
			if a.Target == nil {
				return fmt.Sprintf("%d: %s with kicker", a.Card.Kicker.Cost, a.Card)
			}
			return fmt.Sprintf("%d: %s on %s %s with kicker",
				a.Card.Kicker.Cost, a.Card, a.targetPronoun(p), a.Target)
		}
		if a.Card.IsLand {
			return fmt.Sprintf("%s", a.Card)
		}
		if a.Target == nil {
			return fmt.Sprintf("%d: %s", a.Card.ManaCost, a.Card)
		}
		return fmt.Sprintf("%d: %s on %s %s",
			a.Card.ManaCost, a.Card, a.targetPronoun(p), a.Target)
	case DeclareAttack:
		return "enter attack step"
	case Attack:
		return fmt.Sprintf("attack with %s", a.Card)
	case Block:
		return fmt.Sprintf("%s blocks %s", a.Card, a.Target)
	case UseForMana:
		return fmt.Sprintf("tap %s for mana", a.Card)
	}
	panic("control should not reach here")
}
