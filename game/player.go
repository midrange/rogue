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
		card.Effects = []*Effect{}
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

	c.Effects = []*Effect{}
	for _, aura := range c.Auras {
		p.RemoveFromBoard(aura)
	}
}

// Returns possible actions when we can play a card from hand, including passing.
func (p *Player) PlayActions(allowSorcerySpeed bool, forHuman bool) []*Action {
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
			if card.IsEnchantCreature && mana >= card.ManaCost && card.HasLegalTarget(p.Game) {
				if forHuman {
					answer = append(answer, &Action{
						Type: ChooseTargetAndMana,
						Card: card,
					})
				} else {
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
		// TODO - add player targets - this assumes all instants target creatures for now
		if card.IsInstant && mana >= card.ManaCost && card.HasLegalTarget(p.Game) {
			if forHuman {
				answer = append(answer, &Action{
					Type: ChooseTargetAndMana,
					Card: card,
				})
			} else {
				for _, target := range p.Game.Creatures() {
					answer = append(answer, &Action{
						Type:   Play,
						Card:   card,
						Target: target,
					})
				}
			}
		}
		// TODO - can a card have a 0 kicker, do we need a nullable value here?
		if card.IsInstant && card.KickerCost > 0 && mana >= card.KickerCost && card.HasLegalTarget(p.Game) {
			if !forHuman {
				for _, target := range p.Game.Creatures() {
					answer = append(answer, &Action{
						Type:   PlayWithKicker,
						Card:   card,
						Target: target,
					})
				}
			}
		}
	}
	return answer
}

// Returns possible actions to generate mana.
func (p *Player) ManaActions() []*Action {
	actions := []*Action{}
	for _, card := range p.Board {
		actions = append(actions, card.ManaActions()...)
	}
	return actions
}

// Returns just the pass action,
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

func (p *Player) Play(action *Action, kicker bool) {
	card := action.Card
	newHand := []*Card{}
	for _, c := range p.Hand {
		if c != card {
			newHand = append(newHand, c)
		}
	}
	card.TurnPlayed = p.Game.Turn
	if card.IsLand {
		p.LandPlayedThisTurn++
	} else {
		if kicker {
			p.SpendMana(card.KickerCost)
		} else {
			p.SpendMana(card.ManaCost)
		}
		for _, permanent := range p.Board {
			permanent.RespondToSpell(card)
		}
	}

	if card.IsInstant {
		card.DoEffect(action, kicker)
		// TODO put instants and sorceries in graveyard (or exile)
	} else {
		// TODO allow for kicked creatures
		p.Board = append(p.Board, card)
		if card.IsEnchantCreature {
			action.Target.Auras = append(action.Target.Auras, card)
		}
	}

	p.Hand = newHand
}

func (p *Player) AddMana() {
	p.ColorlessManaPool += 1
}

func (p *Player) Print(position int, hideCards bool, gameWidth int) {
	if position == 0 {
		PrintRowOfCards(p.NonLandPermanents(), gameWidth)
		PrintRowOfCards(p.Lands(), gameWidth)
		PrintRowOfCards(p.Hand, gameWidth)
		fmt.Printf("\n%v", p.AvatarString(position, gameWidth))
	} else {
		fmt.Printf("\n%v\n", p.AvatarString(position, gameWidth))
		PrintRowOfCards(p.Hand, gameWidth)
		PrintRowOfCards(p.Lands(), gameWidth)
		PrintRowOfCards(p.NonLandPermanents(), gameWidth)
	}
}

func (p *Player) Lands() []*Card {
	lands := []*Card{}
	for _, card := range p.Board {
		if card.IsLand {
			lands = append(lands, card)
		}
	}
	return lands
}

func (p *Player) NonLandPermanents() []*Card {
	other := []*Card{}
	for _, card := range p.Board {
		if !card.IsLand && !card.IsEnchantCreature {
			other = append(other, card)
		}
	}
	return other
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
