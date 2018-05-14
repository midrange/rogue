package main

import (
	"bufio"
	"fmt"
	"github.com/midrange/rogue/game"
	"math/rand"
	"os"
	"strconv"
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
	fmt.Println("1) Human vs AI")
	fmt.Println("2) AI vs AI")
	fmt.Print("\nEnter a number: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func playHumanVsComputer() {
	newGame := game.NewGame(game.Stompy(), game.Stompy())
	for {
		actions := newGame.Actions()
		if newGame.Attacker() == newGame.Players[0] {
			if len(actions) == 1 {
				newGame.TakeAction(actions[0])
			} else {
				// get a human move
				newGame.Print()
				promptForAction(newGame, actions)
			}
		} else {
			// get a computer move
			randomAction := actions[rand.Int()%len(actions)]
			fmt.Printf("Computer %v\n", randomAction.String())
			newGame.TakeAction(randomAction)
		}
		if newGame.IsOver() {
			break
		}
	}
}

func promptForAction(newGame *game.Game, actions []*game.Action) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n## Turn %v | %v\n", newGame.Turn, newGame.Phase)
	for index, action := range actions {
		fmt.Printf("%v) %v\n", index, action.String())
	}
	fmt.Print("\nEnter a number: ")
	text, _ := reader.ReadString('\n')
	int_choice, err := strconv.Atoi(strings.TrimSpace(text))
	if err == nil && int_choice >= 0 && int_choice < len(actions) {
		newGame.TakeAction(actions[int_choice])
	} else {
		promptForAction(newGame, actions)
	}
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
