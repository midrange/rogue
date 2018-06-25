package game

import (
	"fmt"
	"github.com/jinzhu/copier"
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

func PlayGame(g *Game, strategy0 Strategy, strategy1 Strategy, printResult bool) PlayerId {
	for !g.IsOver() {
		strategy := []Strategy{strategy0, strategy1}[g.PriorityIndex()]
		action := strategy.Action(g)
		g.TakeAction(action)
	}

	if g.Players[0].Lost() {
		if g.Players[1].Lost() {
			if printResult {
				fmt.Println("The game is a draw.")
			}
		} else {
			if printResult {
				fmt.Printf("%s wins.\n", strategy1)
			}
		}
		return g.DefenderId()
	} else {
		if printResult {
			fmt.Printf("%s wins.\n", strategy0)
		}
		return g.AttackerId()
	}
}

type SimpleMonteCarloBot struct{}

func (b *SimpleMonteCarloBot) String() string {
	return "SimpleMonteCarloBot"
}

func (b *SimpleMonteCarloBot) Action(g *Game) *Action {
	actions := g.Actions(false)

	scores := []int{}
	topScoreIndex := 0
	topScore := 0
	iterations := 1000

	if len(actions) != 1 {
		for moveIndex, _ := range actions {
			score := b.calcWinRate(g, moveIndex, iterations)
			if score > topScore {
				topScore = score
				topScoreIndex = moveIndex
			}
			scores = append(scores, score)
			moveIndex += 1
		}
	}

	return actions[topScoreIndex]
}

func (b *SimpleMonteCarloBot) calcWinRate(g *Game, moveIndex int, iterations int) int {
	wins := 0
	losses := 0
	for i := 0; i < iterations; i++ {
		cloneGame := Game{}
		copier.Copy(&cloneGame, &g)

		move := cloneGame.Actions(false)[moveIndex]
		cloneGame.TakeAction(move)
		winner := PlayGame(&cloneGame, &RandomBot{}, &RandomBot{}, false)
		if winner == g.PriorityId {
			wins += 1
		} else {
			losses += 1
		}
	}
	return wins / (losses + wins*1.0) * 100
}
