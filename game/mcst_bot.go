package game

import (
	"fmt"
	"math"
	"math/rand"
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
		calculationTime: 1.0,
		maxMoves:        10000,
		plays:           map[string]int{},
		wins:            map[string]int{},
	}
	return mcst
}

// An Action, plus the EndState it reaches upon being played.
type ActionState struct {
	EndState []byte
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
	// start := time.Now()
	gameState := g.Serialized()
	actionStates := g.ActionStates(gameState)
	for {
		// print a spinner
		mb.doPlayOut(gameState, g, actionStates)
		games += 1
		//if time.Since(start).Seconds() > mb.calculationTime {
		//	break
		///}
		if games > 10 {
			break
		}
	}
	fmt.Println("Simulated ", games, " games.")

	bestActionState := actionStates[0]
	bestScore := 0
	for _, as := range actionStates {
		endStateStr := string(as.EndState[:])
		if mb.plays[endStateStr] > 0 {
			score := mb.wins[endStateStr] / mb.plays[endStateStr]
			if score > bestScore {
				bestScore = score
				bestActionState = as
			}
		}
		fmt.Printf("%s: %.2f (%d / %d)\n", as.Action.ShowTo(g.Priority()), float64(mb.wins[endStateStr])/float64(mb.plays[endStateStr]), mb.wins[endStateStr], mb.plays[endStateStr])
	}
	return bestActionState.Action
}

func (mb *McstBot) doPlayOut(gameState []byte, g *Game, actionStates []*ActionState) {
	visitedStates := [][]byte{}
	expand := true

	cloneGame := DeserializeGameState(gameState)
	t := 0
	bestActionState := &ActionState{}
	for t = 0; t < mb.maxMoves; t++ {
		fmt.Println(t)
		if t != 0 {
			actionStates = cloneGame.ActionStates(gameState)
		}
		statsForAllPlays := true
		for _, actionState := range actionStates {
			endStateStr := string(actionState.EndState[:])
			if mb.plays[endStateStr] == 0 {
				statsForAllPlays = false
				break
			}

		}

		if statsForAllPlays {
			// decide best play based on prior simulatons
			logTotalPlays := 0.0
			for _, as := range actionStates {
				endStateStr := string(as.EndState[:])
				logTotalPlays += float64(mb.plays[endStateStr])
			}
			logTotalPlays = math.Log(float64(logTotalPlays))

			bestScore := 0.0
			for _, as := range actionStates {
				endStateStr := string(as.EndState[:])
				winRatio := float64(mb.wins[endStateStr]) / float64(mb.plays[endStateStr])
				logPlayRatio := float64(logTotalPlays) / float64(mb.plays[endStateStr])
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
		fmt.Printf("Priority pre action %d taking action %s\n", cloneGame.PriorityId, bestActionState.Action.ShowTo(cloneGame.Priority()))
		cloneGame.TakeAction(bestActionState.Action)
		fmt.Printf("Priority post action %d\n", cloneGame.PriorityId)
		gameState = bestActionState.EndState

		if fmt.Sprintf("%s", cloneGame.Serialized()) != fmt.Sprintf("%s", gameState) {
			fmt.Printf("%s\n", cloneGame.Stack)
			fmt.Printf("%s\n", bestActionState.Action)
			fmt.Printf("%s\n", cloneGame.Serialized())
			fmt.Printf("%s\n", gameState)
			for x := 0; x < len(gameState); x++ {
				if cloneGame.Serialized()[x] != gameState[x] {
					fmt.Printf("%d %v", x, gameState[x])
					break
				}
			}
			panic("these should be same")
		}

		// update stats
		endStateStr := string(gameState[:])
		if expand && mb.plays[endStateStr] == 0 {
			expand = false
		}
		if cloneGame.IsOver() {
			break
		}
		visitedStates = append(visitedStates, gameState)
	}

	for _, es := range visitedStates {
		endStateStr := string(es[:])
		mb.plays[endStateStr] += 1
		if cloneGame.Defender().Lost() && cloneGame.DefenderId() == g.PriorityId {
			mb.wins[endStateStr] += 1
		}
	}
}
