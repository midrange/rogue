package game

import (
	"fmt"
)

type Action struct {
	Type ActionType

	Card          *Card
	Cost          *Cost
	Selected      []*Permanent // for non-targeted effects, such as in Snap
	Source        *Permanent   // for targeted effects
	SpellTarget   *Action
	Target        *Permanent
	With          *Permanent // for attacking
	WithAlternate bool
	WithKicker    bool
	WithPhyrexian bool
}

//go:generate stringer -type=ActionType
type ActionType int

const (
	Pass ActionType = iota
	Play
	DeclareAttack
	Attack
	Block
	UseForMana
	ChooseTargetAndMana
	Activate
	PassPriority
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
				return fmt.Sprintf("%s: %s with kicker", a.Card.Kicker.Cost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s with kicker",
				a.Card.Kicker.Cost, a.Card, a.targetPronoun(p), a.Target)
		}
		if a.Card.IsLand() {
			return fmt.Sprintf("%s", a.Card)
		}
		if a.Target == nil {
			return fmt.Sprintf("%s: %s", a.Card.CastingCost, a.Card)
		}
		return fmt.Sprintf("%s: %s on %s %s",
			a.Card.CastingCost, a.Card, a.targetPronoun(p), a.Target)
	case DeclareAttack:
		return "enter attack step"
	case Attack:
		return fmt.Sprintf("attack with %s", a.With)
	case Block:
		return fmt.Sprintf("%s blocks %s", a.With, a.Target)
	case UseForMana:
		return fmt.Sprintf("tap %s for mana", a.Source)
	case Activate:
		return fmt.Sprintf("use %s", a.Source)
	}
	panic("control should not reach here")
}
