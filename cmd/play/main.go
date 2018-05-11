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
	game := game.NewGame(game.Stompy(), game.Stompy())

	for {
		actions := game.Actions()
		randomAction := actions[rand.Int()%len(actions)]
		fmt.Println(randomAction.Type, randomAction.Card)
		game.TakeAction(randomAction)
		if game.IsOver() {
			break
		}
	}
	fmt.Println("Game over!")
}
