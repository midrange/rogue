package game

import (
	"fmt"
	"math/rand"
)

// The only thing a strategy has to do is to decide an action based on the current
// state of a game.
type Strategy interface {
	// String() should return just the overall name of a strategy, not anything
	// about its state.
	String() string

	Action(g *Game) *Action
}

type RandomBot struct{}

func (b *RandomBot) String() string {
	return "RandomBot"
}

func (b *RandomBot) Action(g *Game) *Action {
	actions := g.Actions(false)
	return actions[rand.Int()%len(actions)]
}

func PlayGame(g *Game, strategy0 Strategy, strategy1 Strategy) {
	for !g.IsOver() {
		strategy := []Strategy{strategy0, strategy1}[g.PriorityIndex()]
		action := strategy.Action(g)
		g.TakeAction(action)
		fmt.Printf("%s %s.\n", strategy, action)
	}

	if g.Players[0].Lost() {
		if g.Players[1].Lost() {
			fmt.Println("The game is a draw.")
		} else {
			fmt.Printf("%s wins.\n", strategy1)
		}
	} else {
		fmt.Printf("%s wins.\n", strategy0)
	}
}
