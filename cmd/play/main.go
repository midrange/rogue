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
		playOutGameRandomly()
	    i++
	}
	fmt.Printf("Played out %v games in 1 second\n", i)
}

func playOutGameRandomly() {
	game := game.NewGame(game.Stompy(), game.Stompy())

	for {
		actions := game.Actions()
		randomAction := actions[rand.Int()%len(actions)]
		game.TakeAction(randomAction)
		if game.IsOver() {
			break
		}
	}
}

