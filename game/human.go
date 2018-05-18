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
	player := game.Priority()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("## Turn %d | %s\n", game.Turn, game.Phase)
		for index, action := range actions {
			fmt.Printf("%d) %s\n", index+1, action.ShowTo(player))
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
	player := game.Priority()
	c := action.Card
	actions := []*Action{}
	for _, target := range game.Creatures() {
		actions = append(actions, &Action{
			Type:   Play,
			Card:   c,
			Target: target,
		})
	}
	mana := game.Priority().AvailableMana()
	if c.IsInstant() && c.Kicker != nil && c.Kicker.CastingCost.Colorless > 0 && mana >= c.Kicker.CastingCost.Colorless {
		for _, target := range game.Creatures() {
			if player.IsLegalTarget(c, target) {
				actions = append(actions, &Action{
					Type:       Play,
					Card:       c,
					WithKicker: true,
					Target:     target,
				})
			}
		}
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		for index, action := range actions {
			fmt.Printf("%d) %s\n", index, action.ShowTo(player))
		}
		fmt.Print("\nEnter a number: ")
		text, _ := reader.ReadString('\n')
		intChoice, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && intChoice >= 0 && intChoice < len(actions) {
			return actions[intChoice]
		}
	}
}
