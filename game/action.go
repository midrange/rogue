package game

import (
	"fmt"
	"strings"
)

type Action struct {
	Type ActionType

	// a faux effect that resolves after a choice-based action, such as returning Scry cards and drawing
	AfterEffect *Effect
	Card        *Card
	// the spell target Card's coming into play effect
	EntersTheBattleFieldSpellTarget *StackObject
	Cost                            *Cost
	// for non-targetted effects, such as in Snap
	Selected []*Permanent
	// whether to switch priority after the action
	ShouldSwitchPriority bool
	// for targeted effects
	Source      *Permanent
	SpellTarget *StackObject
	Target      *Permanent
	// for attacking
	With          *Permanent
	WithAlternate bool
	WithKicker    bool
	WithNinjitsu  bool
	WithPhyrexian bool
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
	forHuman := false
	switch a.Type {
	case PassPriority:
		if len(p.game.Stack) > 0 {
			return fmt.Sprintf("%s", p.game.Stack[len(p.game.Stack)-1])
		}
		if p.game.Phase == Upkeep ||
			p.game.Phase == Draw ||
			p.game.Phase == CombatDamage ||
			p.game.Phase == Main1 ||
			p.game.Phase == Main2 {
			return fmt.Sprintf("end %s", p.game.Phase)
		}
		return "Pass priority"
	case Pass:
		if len(p.game.Stack) > 0 {
			return fmt.Sprintf("resolve %s", p.game.Stack[len(p.game.Stack)-1])
		}
		if p.game.Phase == Upkeep ||
			p.game.Phase == Draw ||
			p.game.Phase == CombatDamage ||
			p.game.Phase == Main1 ||
			p.game.Phase == Main2 {
			return fmt.Sprintf("agree to end %s", p.game.Phase)
		}
		return "Pass"
	case ChooseTargetAndMana:
		forHuman = true
		fallthrough
	case Play:
		if a.WithNinjitsu {
			return fmt.Sprintf("%s: %s", a.Card.Ninjitsu, a.Card)
		}
		if a.WithAlternate {
			if a.Target == nil {
				return fmt.Sprintf("%s: %s", a.Card.AlternateCastingCost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s",
				a.Card.AlternateCastingCost, a.Card, a.targetPronoun(p), a.Target)
		}
		if a.WithPhyrexian {
			if a.Target == nil {
				return fmt.Sprintf("%s: %s", a.Card.PhyrexianCastingCost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s",
				a.Card.PhyrexianCastingCost, a.Card, a.targetPronoun(p), a.Target)
		}
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
			if forHuman && (a.Card.AlternateCastingCost != nil || a.Card.PhyrexianCastingCost != nil) {
				return fmt.Sprintf("%s", a.Card)
			}
			return fmt.Sprintf("%s: %s", a.Card.CastingCost, a.Card)
		}
		return fmt.Sprintf("%s: %s on %s %s",
			a.Card.CastingCost, a.Card, a.targetPronoun(p), a.Target)
	case Attack:
		return fmt.Sprintf("Attack with %s", a.With)
	case Block:
		return fmt.Sprintf("%s blocks %s", a.With, a.Target)
	case UseForMana:
		return fmt.Sprintf("Tap %s for mana", a.Source)
	case Activate:
		return fmt.Sprintf("Use %s", a.Source)
	case MakeChoice:
		if a.AfterEffect.EffectType == ReturnScryCardsDraw {
			topStrings := []string{}
			for _, cn := range a.AfterEffect.ScryCards[0] {
				topStrings = append(topStrings, fmt.Sprintf("%s", cn))
			}
			bottomStrings := []string{}
			for _, cn := range a.AfterEffect.ScryCards[1] {
				bottomStrings = append(bottomStrings, fmt.Sprintf("%s", cn))
			}
			return fmt.Sprintf("Top: %s, Bottom: %s", strings.Join(topStrings, ", "), strings.Join(bottomStrings, ", "))
		}
		if a.AfterEffect.EffectType == ReturnCardsToTopDraw {
			nameStrings := []string{}
			for _, cn := range a.AfterEffect.Cards {
				nameStrings = append(nameStrings, fmt.Sprintf("%s", cn))
			}
			return fmt.Sprintf(strings.Join(nameStrings, ", "))
		}
		if a.AfterEffect.EffectType == ShuffleDraw {
			return fmt.Sprintf("Shuffle")
		}
		if a.AfterEffect.EffectType == DelverScryReveal {
			return fmt.Sprintf("Reveal %s", a.AfterEffect.Cards[0])
		}
		if a.AfterEffect.EffectType == DelverScryNoReveal {
			return fmt.Sprintf("Don't reveal %s", a.AfterEffect.Cards[0])
		}
		return fmt.Sprintf("Choose %s", a.AfterEffect)
	}
	fmt.Println("action is ", a)
	panic("control should not reach here")
}
