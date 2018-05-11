package main

import (
	"fmt"
	"github.com/midrange/rogue/game"
	"math/rand"
	"time"
)

func main() {
	var i int
	for start := time.Now(); time.Since(start) < time.Second; {
		playOutGameRandomly(i)
	    i++
	}
	fmt.Printf("Played out %v games in 1 second\n", i)
}


func playOutGameRandomly(gameNumber int) {
	game := game.NewGame()
	for {
		actions := game.Actions()
		randomAction := actions[rand.Int() % len(actions)]
		game.TakeAction(randomAction)
		if game.IsOver() { 
			break 
		}
	}

}

