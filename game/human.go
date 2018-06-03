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
	allowSorcerySpeed := game.PriorityId == game.AttackerId()
	for {
		reader := bufio.NewReader(os.Stdin)
		whoseTurn := "your turn"
		if player != game.Attacker() {
			whoseTurn = "opponent's turn"
		}
		fmt.Printf("## Turn %d | %s (%s)\n", game.Turn, game.Phase, whoseTurn)
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
				return promptForTargetAndMana(allowSorcerySpeed, game, action)
			}
			return actions[intChoice]
		}
	}
}

func promptForTargetAndMana(allowSorcerySpeed bool, game *Game, action *Action) *Action {
	player := game.Priority()
	card := action.Card
	actions := []*Action{}
	if allowSorcerySpeed {
		actions = player.appendActionsIfNonInstant(actions, card, false)
	}

	if card.IsInstant() && (player.HasLegalPermanentTarget(card) || card.HasCreatureTargets() == false) {
		actions = player.appendActionsForInstant(actions, card)
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
