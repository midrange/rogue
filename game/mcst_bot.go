package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
	MCST is a Strategy that implements a monte carlo search tree.
*/

type McstBot struct {
	// count times states have been reached in playouts
	plays map[string]int
	// count times states that have been reached in playouts led to a win
	wins map[string]int
	// the amount of time to call run_simulation as much as possible
	calculationTime float64
	// the max_moves for any play out
	maxMoves int
	/*
		Larger C encourages more exploration of the possibilities,
		smaller causes the AI to prefer concentrating on known good moves
	*/
	C float64
}

func NewMcstBot() *McstBot {
	mcst := &McstBot{
		C:               1.4,
		calculationTime: 10.0,
		maxMoves:        10000,
		plays:           map[string]int{},
		wins:            map[string]int{},
	}
	return mcst
}

// An Action, plus the EndState it reaches upon being played.
type ActionState struct {
	EndState string
	Action   *Action
}

func (mb *McstBot) String() string {
	return "McstBot"
}

// Return the best play, after simulating possible plays and updating plays and wins stats.
func (mb *McstBot) Action(g *Game) *Action {
	legal := g.Actions(false)
	if len(legal) == 1 {
		return legal[0]
	}
	games := 0
	start := time.Now()
	gameState := g.Serialized()
	for {
		// print a spinner
		mb.doPlayOut(gameState, g)
		games += 1
		if time.Since(start).Seconds() > mb.calculationTime {
			break
		}
	}
	fmt.Println("Simulated ", games, " games.")

	actionStates := g.ActionStates()

	bestActionState := actionStates[0]
	bestScore := 0
	for _, as := range actionStates {
		if mb.plays[as.EndState] > 0 {
			score := mb.wins[as.EndState] / mb.plays[as.EndState]
			if score > bestScore {
				bestScore = score
				bestActionState = as
			}
		}
		fmt.Printf("%s: %.2f (%d / %d)\n", as.Action.ShowTo(g.Priority()), float64(mb.wins[as.EndState])/float64(mb.plays[as.EndState]), mb.wins[as.EndState], mb.plays[as.EndState])
	}
	return bestActionState.Action
}

func (mb *McstBot) doPlayOut(gameState []byte, g *Game) {
	visitedStates := []string{}
	expand := true

	cloneGame := DeserializeGameState(gameState)
	t := 0
	for t = 0; t < mb.maxMoves; t++ {
		actionStates := cloneGame.ActionStates()
		statsForAllPlays := true
		for _, actionState := range actionStates {
			if mb.plays[actionState.EndState] == 0 {
				statsForAllPlays = false
				break
			}

		}

		bestActionState := actionStates[0]
		if statsForAllPlays {
			// decide best play based on prior simulatons
			logTotalPlays := 0.0
			for _, as := range actionStates {
				logTotalPlays += float64(mb.plays[as.EndState])
			}
			logTotalPlays = math.Log(float64(logTotalPlays))

			bestScore := 0.0
			for _, as := range actionStates {
				winRatio := float64(mb.wins[as.EndState]) / float64(mb.plays[as.EndState])
				logPlayRatio := float64(logTotalPlays) / float64(mb.plays[as.EndState])
				score := winRatio + mb.C*math.Sqrt(logPlayRatio)
				if score > bestScore {
					bestScore = score
					bestActionState = as
				}
			}
		} else {
			// otherwise play randomly
			bestActionState = actionStates[rand.Intn(len(actionStates))]
		}

		cloneGame.TakeAction(bestActionState.Action)

		// update stats
		if expand && mb.plays[bestActionState.EndState] == 0 {
			expand = false
		}
		if cloneGame.IsOver() {
			break
		}
		visitedStates = append(visitedStates, bestActionState.EndState)
	}

	for _, es := range visitedStates {
		mb.plays[es] += 1
		if cloneGame.Defender().Lost() && cloneGame.DefenderId() == g.PriorityId {
			mb.wins[es] += 1
		}
	}
}
