package game

import (
	"fmt"
)

type Player struct {
	Life               int
	Hand               []*Card
	Board              []*Card
	Opponent           *Player
	Deck               *Deck
	LandPlayedThisTurn int
}

func NewPlayer(deck *Deck) *Player {
	hand := []*Card{}
	for i := 0; i < 7; i++ {
		hand = append(hand, deck.Draw())
	}
	return &Player{
		Life:  20,
		Hand:  hand,
		Deck:  deck,
		Board: []*Card{},
	}
}

func (p *Player) AvailableMana() int {
	answer := 0
	for _, card := range p.Board {
		if card.IsLand {
			answer += 1
		}
	}
	return answer
}

func (p *Player) Draw() {
	card := p.Deck.Draw()
	if card != nil {
		p.Hand = append(p.Hand, card)
	}
}

func (p *Player) Untap() {
	p.LandPlayedThisTurn = 0
	for _, card := range p.Board {
		card.Tapped = false
	}
}

func (p *Player) Lost() bool {
	return p.Life <= 0 || p.Deck.FailedToDraw
}

// Automatically spends the given amount of mana.
// Panics if we do not have that much.
func (p *Player) SpendMana(amount int) {
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
	}
}

// Possible actions when we can play a card from hand, including passing.
func (p *Player) PlayActions(allowSorcerySpeed bool) []*Action {
	cardNames := make(map[CardName]bool)
	answer := []*Action{&Action{Type: Pass}}
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
		}
	}
	return answer
}

// Possible actions when we are announcing attacks, including passing.
func (p *Player) AttackActions() []*Action {
	answer := []*Action{&Action{Type: Pass}}
	for _, card := range p.Board {
		if card.IsCreature && !card.Attacking {
			answer = append(answer, &Action{Type: Attack, Card: card})
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
	if card.IsLand {
		p.LandPlayedThisTurn++
	}
	if card.IsCreature {
		p.SpendMana(card.ManaCost)
	}
	p.Board = append(p.Board, card)
}

func (p *Player) Print(position int, hideCards bool) {
	if position == 0 {
		for _, card := range p.Board {
			card.Print()
		}
		fmt.Printf("\n")
		for _, card := range p.Hand {
			card.Print()
		}
		fmt.Printf("\n")
		fmt.Println("Player ", position, " Life: ", p.Life)
	} else {
		fmt.Println("Player ", position, " Life: ", p.Life)
		for _, card := range p.Hand {
			card.Print()
		}
		fmt.Printf("\n")
		for _, card := range p.Board {
			card.Print()
		}
		fmt.Printf("\n")
	}

}