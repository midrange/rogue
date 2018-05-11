package game

import ()

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
		card.Blocking = nil
		card.DamageOrder = []*Card{}
	}
}

func (p *Player) EndTurn() {
	for _, card := range p.Board {
		card.Damage = 0
	}
}

func (p *Player) RemoveFromBoard(c *Card) {
	newBoard := []*Card{}
	for _, card := range p.Board {
		if card != c {
			newBoard = append(newBoard, card)
		}
	}
	p.Board = newBoard
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

func (p *Player) BlockActions() []*Action {
	answer := []*Action{&Action{Type: Pass}}
	attackers := []*Card{}
	for _, card := range p.Opponent.Board {
		if card.Attacking {
			attackers = append(attackers, card)
		}
	}
	for _, card := range p.Board {
		if card.Blocking == nil && !card.Tapped {
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
	if card.IsLand {
		p.LandPlayedThisTurn++
	}
	if card.IsCreature {
		p.SpendMana(card.ManaCost)
	}
	p.Board = append(p.Board, card)
}
