package game

import (
	"fmt"
	"log"
)

type Player struct {
	Life               int
	ColorlessManaPool  int
	Hand               []*Card
	Board              []*Card
	Opponent           *Player
	Game               *Game
	Deck               *Deck
	LandPlayedThisTurn int
}

// The caller should set Game and Opponent
func NewPlayer(deck *Deck) *Player {
	p := &Player{
		Life:  20,
		Hand:  []*Card{},
		Board: []*Card{},
		Deck:  deck,
	}
	for i := 0; i < 7; i++ {
		p.Draw()
	}
	return p
}

func (p *Player) Draw() {
	card := p.Deck.Draw()
	p.AddToHand(card)
}

func (p *Player) AddToHand(c *Card) {
	// a card would be nil here if you attempted to start a game with less than 7 cards in your deck
	if c == nil {
		return
	}
	c.Owner = p
	p.Hand = append(p.Hand, c)
}

func (p *Player) AvailableMana() int {
	answer := 0
	for _, card := range p.Board {
		if card.IsLand && !card.Tapped {
			answer += 1
		}
	}
	answer += p.ColorlessManaPool
	return answer
}

func (p *Player) Untap() {
	p.LandPlayedThisTurn = 0
	for _, card := range p.Board {
		card.RespondToUntapPhase()
	}
}

func (p *Player) Lost() bool {
	return p.Life <= 0 || p.Deck.FailedToDraw
}

// Automatically spends the given amount of mana.
// Panics if we do not have that much.
func (p *Player) SpendMana(amount int) {
	if p.ColorlessManaPool >= amount {
		p.ColorlessManaPool -= amount
		return
	} else {
		amount -= p.ColorlessManaPool
		p.ColorlessManaPool = 0

	}
	for _, card := range p.Board {
		if amount == 0 {
			return
		}
		if card.IsLand && !card.Tapped {
			card.Tapped = true
			amount -= 1
		}
	}
	if amount > 0 {
		panic("could not spend mana")
	}
}

func (p *Player) EndCombat() {
	for _, card := range p.Board {
		card.Attacking = false
		card.Blocking = nil
		card.DamageOrder = []*Card{}
	}
}

func (p *Player) EndPhase() {
	p.ColorlessManaPool = 0
}

func (p *Player) EndTurn() {
	for _, card := range p.Board {
		card.Damage = 0
	}
	p.LandPlayedThisTurn = 0
	p.EndPhase()
}

func (p *Player) Creatures() []*Card {
	answer := []*Card{}
	for _, card := range p.Board {
		if card.IsCreature {
			answer = append(answer, card)
		}
	}
	return answer
}

func (p *Player) RemoveFromBoard(c *Card) {
	newBoard := []*Card{}
	for _, card := range p.Board {
		if card != c {
			newBoard = append(newBoard, card)
		}
	}
	p.Board = newBoard

	if c.Name == Rancor {
		p.AddToHand(NewCard(Rancor))
	} else {
		// If we had a graveyard we would put the card in the graveyard here
	}

	for _, aura := range c.Auras {
		p.RemoveFromBoard(aura)
	}
}

// Possible actions when we can play a card from hand, including passing.
func (p *Player) PlayActions(allowSorcerySpeed bool) []*Action {
	cardNames := make(map[CardName]bool)
	answer := []*Action{}
	if allowSorcerySpeed {
		answer = append(answer, &Action{Type: PassTurn})
	} else {
		answer = append(answer, &Action{Type: PassPriority})
	}

	mana := p.AvailableMana()
	for _, card := range p.Hand {
		// Don't re-check playing duplicate cards
		if cardNames[card.Name] {
			continue
		}
		cardNames[card.Name] = true

		if allowSorcerySpeed {
			if card.IsLand && p.LandPlayedThisTurn == 0 {
				answer = append(answer, &Action{Type: Play, Card: card})
			}
			if card.IsCreature && mana >= card.ManaCost {
				answer = append(answer, &Action{Type: Play, Card: card})
			}
			if card.IsEnchantCreature && mana >= card.ManaCost {
				for _, target := range p.Game.Creatures() {
					answer = append(answer, &Action{
						Type:   Play,
						Card:   card,
						Target: target,
					})
				}
			}
		}
	}
	return answer
}

// Possible actions to generate mana.
func (p *Player) ManaActions() []*Action {
	actions := []*Action{}
	for _, card := range p.Board {
		actions = append(actions, card.ManaActions()...)
	}
	return actions
}

// Possible actions when we are announcing attacks, including passing.
func (p *Player) PassAction() *Action {
	return &Action{Type: PassPriority}
}

// Returns the possible actions of type 'Attack'.
func (p *Player) AttackActions() []*Action {
	if p.Game.Phase != DeclareAttackers {
		log.Fatalf("do not call AttackActions in phase %s", p.Game.Phase)
	}
	answer := []*Action{}
	for _, card := range p.Board {
		if card.IsCreature && !card.Attacking && !card.Tapped && card.TurnPlayed != p.Game.Turn {
			answer = append(answer, &Action{Type: Attack, Card: card})
		}
	}
	return answer
}

// Returns the possible actions of type 'Block'.
func (p *Player) BlockActions() []*Action {
	answer := []*Action{}
	attackers := []*Card{}
	for _, card := range p.Opponent.Board {
		if card.Attacking {
			attackers = append(attackers, card)
		}
	}
	for _, card := range p.Board {
		if card.Blocking == nil && !card.Tapped && card.IsCreature {
			for _, attacker := range attackers {
				answer = append(answer, &Action{
					Type:   Block,
					Card:   card,
					Target: attacker,
				})
			}
		}
	}
	return answer
}

func (p *Player) Play(card *Card) {
	newHand := []*Card{}
	for _, c := range p.Hand {
		if c != card {
			newHand = append(newHand, c)
		}
	}
	card.TurnPlayed = p.Game.Turn
	if card.IsLand {
		p.LandPlayedThisTurn++
	}
	if card.IsCreature {
		p.SpendMana(card.ManaCost)
	}

	for _, permanent := range p.Board {
		if !card.IsLand {
			permanent.RespondToSpell(card)
		}
	}

	p.Board = append(p.Board, card)
	p.Hand = newHand
}

func (p *Player) AddMana() {
	p.ColorlessManaPool += 1
}

func (p *Player) Print(position int, hideCards bool, gameWidth int) {
	if position == 0 {
		PrintRowOfCards(p.Board, gameWidth)
		PrintRowOfCards(p.Hand, gameWidth)
		fmt.Printf("\n%v", p.AvatarString(position, gameWidth))
	} else {
		fmt.Printf("\n%v\n", p.AvatarString(position, gameWidth))
		PrintRowOfCards(p.Hand, gameWidth)
		PrintRowOfCards(p.Board, gameWidth)
	}
}

func (p *Player) AvatarString(position int, gameWidth int) string {
	playerString := ""
	for x := 0; x < (gameWidth-len(playerString))/2; x++ {
		playerString += " "
	}
	playerString += fmt.Sprintf("<Life: %v> Player %v <Mana: %v>", p.Life, position, p.ColorlessManaPool)
	return playerString
}

func PrintRowOfCards(cards []*Card, gameWidth int) {
	asciiImages := [][CARD_HEIGHT][CARD_WIDTH]string{}
	for _, card := range cards {
		asciiImages = append(asciiImages, card.AsciiImage(false))
	}
	for row := 0; row < CARD_HEIGHT; row++ {
		for x := 0; x < (gameWidth-len(cards)*(CARD_WIDTH+1))/2; x++ {
			fmt.Printf(" ")
		}
		for _, bitmap := range asciiImages {
			for _, char := range bitmap[row] {
				fmt.Printf(char)
			}
			fmt.Printf(" ")
		}
		fmt.Printf("%v", "\n")
	}
}

// GetCreature gets the first creature in play with the given name.
// It returns nil if there is no such creature in play.
func (p *Player) GetCreature(name CardName) *Card {
	for _, card := range p.Board {
		if card.Name == name {
			return card
		}
	}
	return nil
}
