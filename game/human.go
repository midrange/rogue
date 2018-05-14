package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Human is a Strategy that uses the terminal to ask a human what move to make.
type Human struct{}

func (h *Human) String() string {
	return "Human"
}

func (h *Human) Action(g *Game) *Action {
	actions := g.Actions(true)
	if len(actions) == 1 {
		return actions[0]
	}

	// Get a human move
	g.Print()
	return promptForAction(g, actions)
}

func promptForAction(game *Game, actions []*Action) *Action {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("## Phase: %s\n", game.Phase)
		for index, action := range actions {
			fmt.Printf("%d) %s\n", index, action)
		}
		fmt.Print("\nEnter a number: ")
		text, _ := reader.ReadString('\n')
		intChoice, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && intChoice >= 0 && intChoice < len(actions) {
			action := actions[intChoice]
			if action.Type == ChooseTargetAndMana {
				return promptForTargetAndMana(game, action)
			}
			return actions[intChoice]
		}
	}
}

func promptForTargetAndMana(game *Game, action *Action) *Action {
	actions := []*Action{}
	for _, target := range game.Creatures() {
		actions = append(actions, &Action{
			Type:   Play,
			Card:   action.Card,
			Target: target,
		})
	}
	mana := action.Card.Owner.AvailableMana()
	if action.Card.KickerCost > 0 && action.Card.KickerCost <= mana {
		for _, target := range game.Creatures() {
			actions = append(actions, &Action{
				Type:   PlayWithKicker,
				Card:   action.Card,
				Target: target,
			})
		}
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		for index, action := range actions {
			fmt.Printf("%d) %s\n", index, action)
		}
		fmt.Print("\nEnter a number: ")
		text, _ := reader.ReadString('\n')
		intChoice, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && intChoice >= 0 && intChoice < len(actions) {
			return actions[intChoice]
		}
	}
}
