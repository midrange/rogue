package main

import (
	"bufio"
	"fmt"
	"github.com/midrange/rogue/game"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	text := showWelcomePrompt()
	if text == "1" {
		playHumanVsComputer()
	} else if text == "2" {
		playComputerVsComputer()
	} else {
		main()
	}
}

func showWelcomePrompt() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n ~~~~~~ Welcome to Rogue ~~~~~~\n")
	fmt.Println("1) Play")
	fmt.Println("2) Watch Computers Play")
	fmt.Print("\nEnter a number: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func playHumanVsComputer() {

}

func playComputerVsComputer() {
	i := 0
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
