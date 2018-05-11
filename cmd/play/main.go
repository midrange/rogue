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
		fmt.Println("game started")
		playOutGameRandomly()
		fmt.Println("game ended")
	    i++
	}

	fmt.Println("played ", i, " games")

}


func playOutGameRandomly() {
	game := game.NewGame()
	fmt.Println("game created")
	for {
		// fmt.Printf("M")
		actions := game.Actions()
		randomAction := actions[rand.Int() % len(actions)]
		game.TakeAction(randomAction)
		if game.IsOver() { break }
	}

}

