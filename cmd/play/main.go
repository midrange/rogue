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
		playOutGameRandomly(i)
		fmt.Println("game ended")
	    i++
	}

	fmt.Println("played ", i, " games")

}


func playOutGameRandomly(gameNumber int) {
	game := game.NewGame()
	fmt.Println("game created")
	var i int
	for {
		// fmt.Printf("M")
		actions := game.Actions()
		randomAction := actions[rand.Int() % len(actions)]
		game.TakeAction(randomAction)
		if game.IsOver() { break }
		i++
		if i % 9999 == 0 {
			fmt.Println(randomAction)
		}
		if i % 10000 == 0 {
			fmt.Println("Game Number ", gameNumber, " Move # ", i)
			game.Print()
		}
	}

}

