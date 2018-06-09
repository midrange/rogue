package game

import ()

/*
	AttackBot is a Strategy that always plays lands, cast creatures,
	then spells, then attacks wiht everyone, and never blocks.
*/

type AttackBot struct{}

func (ab *AttackBot) String() string {
	return "AttackBot"
}

func (ab *AttackBot) Action(g *Game) *Action {
	actions := g.Actions(false)
	if len(actions) == 1 {
		return actions[0]
	}
	bestAction := actions[0]

	for _, a := range actions {
		if a.Type == Play {
			if bestAction.Type == Play {
				if a.Card.CastingCost != nil && a.Card.CastingCost.Colorless > bestAction.Card.CastingCost.Colorless {
					bestAction = a
				}
			} else {
				bestAction = a
			}
		}
	}
	for _, a := range actions {
		if a.Type == Play && a.Card.IsCreature() {
			if a.Card.CastingCost != nil && a.Card.CastingCost.Colorless > bestAction.Card.CastingCost.Colorless {
				bestAction = a
			}
		}
	}
	for _, a := range actions {
		if a.Type == Play && a.Card.IsLand() {
			bestAction = a
		}
	}
	for _, a := range actions {
		if a.Type == Attack {
			bestAction = a
		}
	}
	return bestAction
}
