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
	EntersTheBattleFieldSpellTarget StackObjectId
	Cost                            *Cost
	// for non-targetted effects, such as in Snap
	Selected []PermanentId
	// whether to switch priority after the action
	ShouldSwitchPriority bool
	// for targeted effects
	Source      PermanentId
	SpellTarget StackObjectId
	Target      PermanentId
	// for attacking
	With          PermanentId
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
	if p.game.Permanent(a.Target).Owner == p.Id {
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
			return fmt.Sprintf("%s", p.game.StackObject(p.game.Stack[len(p.game.Stack)-1]))
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
			return fmt.Sprintf("resolve %s", p.game.StackObject(p.game.Stack[len(p.game.Stack)-1]))
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
			if a.Target == NoPermanentId {
				return fmt.Sprintf("%s: %s", a.Card.AlternateCastingCost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s",
				a.Card.AlternateCastingCost, a.Card, a.targetPronoun(p), a.Target)
		}
		if a.WithPhyrexian {
			if a.Target == NoPermanentId {
				return fmt.Sprintf("%s: %s", a.Card.PhyrexianCastingCost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s",
				a.Card.PhyrexianCastingCost, a.Card, a.targetPronoun(p), p.game.Permanent(a.Target))
		}
		if a.WithKicker {
			if a.Target == NoPermanentId {
				return fmt.Sprintf("%s: %s with kicker", a.Card.Kicker.Cost, a.Card)
			}
			return fmt.Sprintf("%s: %s on %s %s with kicker",
				a.Card.Kicker.Cost, a.Card, a.targetPronoun(p), p.game.Permanent(a.Target))
		}
		if a.Card.IsLand() {
			return fmt.Sprintf("%s", a.Card)
		}
		if a.Target == NoPermanentId {
			if forHuman && (a.Card.AlternateCastingCost != nil || a.Card.PhyrexianCastingCost != nil) {
				return fmt.Sprintf("%s", a.Card)
			}
			return fmt.Sprintf("%s: %s", a.Card.CastingCost, a.Card)
		}
		if len(a.Selected) > 0 {
			cardNames := []string{}
			for _, perm := range a.Selected {
				cardNames = append(cardNames, fmt.Sprintf("%s", p.game.Permanent(perm).Card.Name))
			}
			return fmt.Sprintf("%s: %s on %s %s (%s)",
				a.Card.CastingCost, a.Card, a.targetPronoun(p), p.game.Permanent(a.Target), strings.Join(cardNames, ", "))
		}
		return fmt.Sprintf("%s: %s on %s %s",
			a.Card.CastingCost, a.Card, a.targetPronoun(p), p.game.Permanent(a.Target))
	case Attack:
		return fmt.Sprintf("Attack with %s", p.game.Permanent(a.With))
	case Block:
		return fmt.Sprintf("%s blocks %s", a.With, p.game.Permanent(a.Target))
	case UseForMana:
		return fmt.Sprintf("Tap %s for mana", p.game.Permanent(a.Source))
	case Activate:
		return fmt.Sprintf("Use %s", p.game.Permanent(a.Source))
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
