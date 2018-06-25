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
	// cache availble actions for states
	stateToActions map[string][]*ActionState
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
		C:               3,
		calculationTime: 3.0,
		maxMoves:        10000,
		plays:           map[string]int{},
		wins:            map[string]int{},
		stateToActions:  map[string][]*ActionState{},
	}
	return mcst
}

// An Action, plus the EndState it reaches upon being played.
type ActionState struct {
	Game   *Game
	Action *Action
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
	for {
		// print a spinner
		mb.doPlayOut(g)
		games++
		if time.Since(start).Seconds() > mb.calculationTime {
			break
		}
	}
	fmt.Println("Simulated ", games, " games.")

	actionStates := g.ActionStates()
	bestActionState := actionStates[0]
	bestScore := 0.0
	for index, as := range actionStates {
		endStateStr := fmt.Sprintf("%s%d", string(as.Game.Serialize()), index)
		if mb.plays[endStateStr] > 0 {
			score := float64(mb.wins[endStateStr]) / float64(mb.plays[endStateStr])
			if score >= bestScore {
				bestScore = score
				bestActionState = as
			}
		}
		fmt.Printf("%s: %.2f (%d / %d)\n", as.Action.ShowTo(g.Priority()), float64(mb.wins[endStateStr])/float64(mb.plays[endStateStr]), mb.wins[endStateStr], mb.plays[endStateStr])
	}

	return bestActionState.Action
}

func (mb *McstBot) doPlayOut(g *Game) {
	visitedStates := []string{}

	cloneGame := CopyGame(g)

	t := 0
	bestActionState := &ActionState{}
	unreachedState := false
	for t = 0; t < mb.maxMoves; t++ {
		currentStateStr := string(cloneGame.Serialize())
		firstActionStr := fmt.Sprintf("%s%d", currentStateStr, 0)
		actionStates := []*ActionState{}
		actionIndex := 0
		if mb.plays[firstActionStr] == 0 {
			actionStates = cloneGame.ActionStates()
			mb.stateToActions[currentStateStr] = actionStates
			unreachedState = true
		} else {
			actionStates = mb.stateToActions[currentStateStr]
		}

		if unreachedState {
			actionIndex = rand.Intn(len(actionStates))
			bestActionState = actionStates[actionIndex]
		} else {
			// decide best play based on prior simulations
			logTotalPlays := 0.0
			endStates := []string{}
			for index, _ := range actionStates {
				endStateStr := fmt.Sprintf("%s%d", currentStateStr, index)
				logTotalPlays += float64(mb.plays[endStateStr])
				endStates = append(endStates, endStateStr)
			}
			logTotalPlays = math.Log(float64(logTotalPlays))

			bestScore := 0.0
			bestActionState = actionStates[0]
			for index, as := range actionStates {
				endStateStr := endStates[index]
				winRatio := float64(mb.wins[endStateStr]) / float64(mb.plays[endStateStr])
				logPlayRatio := float64(logTotalPlays) / float64(mb.plays[endStateStr])
				score := winRatio + mb.C*math.Sqrt(logPlayRatio)
				if score >= bestScore {
					bestScore = score
					bestActionState = as
					actionIndex = index
				}
			}
		}

		cloneGame.TakeAction(bestActionState.Action)

		if cloneGame.IsOver() {
			break
		}
		visitedStates = append(visitedStates, fmt.Sprintf("%s%d", currentStateStr, actionIndex))
	}

	for _, es := range visitedStates {
		mb.plays[es] += 1
		if cloneGame.Players[g.DefenderId()].Lost() {
			mb.wins[es] += 1
		}
	}
}
