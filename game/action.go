package game

import (
	"fmt"
)

type Action struct {
	Type ActionType

	// a faux effect that resolves after a choice-based action, such as returning Scry cards and drawing
	AfterEffect                     *Effect
	Card                            *Card
	EntersTheBattleFieldSpellTarget *StackObject // the spell target Card's coming into play effect
	Cost                            *Cost
	Selected                        []*Permanent // for non-targetted effects, such as in Snap
	ShouldSwitchPriority            bool         // whether to switch priority after the action
	Source                          *Permanent   // for targeted effects
	SpellTarget                     *StackObject
	Target                          *Permanent
	With                            *Permanent // for attacking
	WithAlternate                   bool
	WithKicker                      bool
	WithNinjitsu                    bool
	WithPhyrexian                   bool
}

//go:generate stringer -type=ActionType
type ActionType int

const (
	Pass ActionType = iota
	Play
	Activate
	Attack
	Block
	ChooseTargetAndMana
	DeclareAttack
	DecideOnChoice
	DeclineChoice
	EntersTheBattlefieldEffect
	MakeChoice
	PassPriority
	UseForMana
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
	case PassPriority:
		return "pass priority"
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
	case MakeChoice:
		if a.AfterEffect.EffectType == ReturnScryCardsDraw {
			return fmt.Sprintf("%s, Top: %s, Bottom: %s", a.AfterEffect.EffectType, a.AfterEffect.ScryCards[0], a.AfterEffect.ScryCards[1])
		}
		if a.AfterEffect.EffectType == ReturnCardsToTopDraw {
			return fmt.Sprintf("%s, %s", a.AfterEffect.EffectType, a.AfterEffect.Cards)
		}
		return fmt.Sprintf("Choose %s", a.AfterEffect)
	}
	fmt.Println("action is ", a)
	panic("control should not reach here")
}
