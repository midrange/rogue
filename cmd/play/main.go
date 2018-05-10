package main

import (
	"fmt"
	"github.com/midrange/rogue/game"
	"math/rand"
)

func main() {
	playOutGameRandomly()
}


func playOutGameRandomly() {
	fmt.Println("Playing out a game between two bots that move randomly.")
	game := game.NewGame()

	for {
		actions := game.Actions()
		randomAction := actions[rand.Int() % len(actions)]
		game.TakeAction(randomAction)
		game.Print()
		if game.IsOver() { break }
	}
	fmt.Println("Game over!")
}
