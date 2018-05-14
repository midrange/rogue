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
		fmt.Printf("## Turn %v | %s\n", game.Turn, game.Phase)
		for index, action := range actions {
			fmt.Printf("%d) %s\n", index+1, action)
		}
		fmt.Print("\nEnter a number: ")
		text, _ := reader.ReadString('\n')
		intChoice, err := strconv.Atoi(strings.TrimSpace(text))
		intChoice--
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
	if action.Card.IsInstant && action.Card.HasKicker && action.Card.Kicker.Cost > 0 && mana >= action.Card.Kicker.Cost && action.Card.HasLegalTarget(action.Card.Owner.Game) {
		for _, target := range action.Card.Owner.Game.Creatures() {
			if target.Targetable(action.Card) {
				actions = append(actions, &Action{
					Type:       Play,
					Card:       action.Card,
					WithKicker: true,
					Target:     target,
				})
			}
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
