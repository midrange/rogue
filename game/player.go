package game

import (
	"fmt"
	"log"
)

type Player struct {
	Board              []*Permanent
	ColorlessManaPool  int
	CreatureDied       bool
	DamageThisTurn     int
	Deck               *Deck
	Hand               []*Card
	Id                 PlayerId
	LandPlayedThisTurn int
	Life               int

	// game should not be included when the player is serialized.
	game *Game
}

// The caller should set game after construction.
func NewPlayer(deck *Deck, id PlayerId) *Player {
	p := &Player{
		Life:  20,
		Hand:  []*Card{},
		Id:    id,
		Board: []*Permanent{},
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
	if c == nil {
		return
	}
	p.Hand = append(p.Hand, c)
}

func (p *Player) AvailableMana() int {
	answer := 0
	for _, card := range p.Board {
		if card.IsLand && !card.Tapped {
			answer += card.Colorless
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

func (p *Player) Opponent() *Player {
	return p.game.Player(p.Id.OpponentId())
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
		card.DamageOrder = []*Permanent{}
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
	p.DamageThisTurn = 0
	p.CreatureDied = false
	p.EndPhase()
}

func (p *Player) Creatures() []*Permanent {
	answer := []*Permanent{}
	for _, card := range p.Board {
		if card.IsCreature {
			answer = append(answer, card)
		}
	}
	return answer
}

func (p *Player) RemoveFromBoard(perm *Permanent) {
	newBoard := []*Permanent{}
	for _, permanent := range p.Board {
		if permanent != perm {
			newBoard = append(newBoard, permanent)
		}
	}
	p.Board = newBoard

	if perm.Name == Rancor {
		p.AddToHand(NewCard(Rancor))
	} else {
		// TODO: make sure it is actually a creature that died
		p.CreatureDied = true
	}

	perm.Effects = []*Effect{}
	for _, aura := range perm.Auras {
		p.RemoveFromBoard(aura)
	}
}

// Returns possible actions when we can play a card from hand, including passing.
func (p *Player) PlayActions(allowSorcerySpeed bool, forHuman bool) []*Action {
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
			if card.IsCreature && mana >= card.CastingCost.Colorless {
				answer = append(answer, &Action{Type: Play, Card: card})
			}
			if card.IsEnchantCreature && mana >= card.CastingCost.Colorless && p.HasLegalTarget(card) {
				if forHuman {
					answer = append(answer, &Action{
						Type: ChooseTargetAndMana,
						Card: card,
					})
				} else {
					for _, target := range p.game.Creatures() {
						answer = append(answer, &Action{
							Type:   Play,
							Card:   card,
							Target: target,
						})
					}
				}
			}
			if card.IsCreature || card.IsEnchantCreature {
				if card.PhyrexianCastingCost != nil && mana >= card.PhyrexianCastingCost.Colorless && p.Life >= card.PhyrexianCastingCost.Life {
					answer = append(answer, &Action{Type: Play, Card: card, WithPhyrexian: true})
				}
			}
		}

		// TODO - add player targets - this assumes all instants target creatures for now
		if card.IsInstant && mana >= card.CastingCost.Colorless && p.HasLegalTarget(card) {
			if forHuman {
				answer = append(answer, &Action{
					Type: ChooseTargetAndMana,
					Card: card,
				})
			} else {
				for _, target := range p.game.Creatures() {
					if p.IsLegalTarget(card, target) {
						answer = append(answer, &Action{
							Type:   Play,
							Card:   card,
							Target: target,
						})
					}
				}
			}
		}

		if card.IsInstant && card.Kicker != nil && card.Kicker.CastingCost.Colorless > 0 && mana >= card.Kicker.CastingCost.Colorless && p.HasLegalTarget(card) {
			if !forHuman {
				for _, target := range p.game.Creatures() {
					if p.IsLegalTarget(card, target) {
						answer = append(answer, &Action{
							Type:       Play,
							Card:       card,
							WithKicker: true,
							Target:     target,
						})
					}
				}
			}
		}

		if card.IsInstant {
			if card.PhyrexianCastingCost != nil && mana >= card.PhyrexianCastingCost.Colorless && p.Life >= card.PhyrexianCastingCost.Life {
				answer = append(answer, &Action{Type: Play, Card: card, WithPhyrexian: true})
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
	return &Action{Type: Pass}
}

// Returns the possible actions of type 'Attack'.
func (p *Player) AttackActions() []*Action {
	if p.game.Phase != DeclareAttackers {
		log.Fatalf("do not call AttackActions in phase %s", p.game.Phase)
	}
	answer := []*Action{}
	for _, card := range p.Board {
		if card.IsCreature && !card.Attacking && !card.Tapped && card.TurnPlayed != p.game.Turn {
			answer = append(answer, &Action{Type: Attack, With: card})
		}
	}
	return answer
}

// Returns the possible actions of type 'Block'.
func (p *Player) BlockActions() []*Action {
	answer := []*Action{}
	attackers := []*Permanent{}
	for _, perm := range p.Opponent().Board {
		if perm.Attacking {
			attackers = append(attackers, perm)
		}
	}
	for _, perm := range p.Board {
		if perm.Blocking == nil && !perm.Tapped && perm.IsCreature {
			for _, attacker := range attackers {
				if perm.CanBlock(attacker) {
					answer = append(answer, &Action{
						Type:   Block,
						With:   perm,
						Target: attacker,
					})
				}
			}
		}
	}
	return answer
}

func (p *Player) Play(action *Action) {
	card := action.Card
	newHand := []*Card{}
	for _, c := range p.Hand {
		if c != card {
			newHand = append(newHand, c)
		}
	}
	p.Hand = newHand

	if card.IsCreature || card.IsInstant {
		if action.WithKicker {
			p.SpendMana(card.Kicker.CastingCost.Colorless)
		} else if action.WithPhyrexian {
			p.SpendMana(card.PhyrexianCastingCost.Colorless)
			p.Life -= card.PhyrexianCastingCost.Life
		} else {
			p.SpendMana(card.CastingCost.Colorless)
		}
		for _, permanent := range p.Board {
			permanent.RespondToSpell()
		}
	}

	if card.IsInstant {
		p.castInstant(card, action.Target, action)
		// TODO put instants and sorceries in graveyard (or exile)
		return
	}

	// Non-instant cards turn into a permanent
	perm := &Permanent{
		Card:  card,
		Owner: p,
	}
	perm.TurnPlayed = p.game.Turn

	if card.IsLand {
		p.LandPlayedThisTurn++
	}

	p.Board = append(p.Board, perm)
	perm.HandleComingIntoPlay()

	if card.IsEnchantCreature {
		action.Target.Auras = append(action.Target.Auras, perm)
	}
}

func (p *Player) castInstant(c *Card, target *Permanent, a *Action) {
	if c.AddsTemporaryEffect {
		target.Effects = append(target.Effects, NewEffect(a))
	}

	if c.Effect != nil {
		target.Plus1Plus1Counters += c.Effect.Plus1Plus1Counters
	}
	if c.Morbid != nil && (p.CreatureDied || p.Opponent().CreatureDied) {
		target.Plus1Plus1Counters += c.Morbid.Plus1Plus1Counters
	}
}

func (p *Player) AddMana(colorless int) {
	p.ColorlessManaPool += colorless
}

func (p *Player) Print(position int, hideCards bool, gameWidth int) {
	if position == 0 {
		PrintRowOfPermanents(p.NonLandPermanents(), gameWidth)
		PrintRowOfPermanents(p.Lands(), gameWidth)
		PrintRowOfCards(p.Hand, gameWidth)
		fmt.Printf("\n%s", p.AvatarString(position, gameWidth))
	} else {
		fmt.Printf("\n%s\n", p.AvatarString(position, gameWidth))
		PrintRowOfCards(p.Hand, gameWidth)
		PrintRowOfPermanents(p.Lands(), gameWidth)
		PrintRowOfPermanents(p.NonLandPermanents(), gameWidth)
	}
}

func (p *Player) Lands() []*Permanent {
	lands := []*Permanent{}
	for _, perm := range p.Board {
		if perm.IsLand {
			lands = append(lands, perm)
		}
	}
	return lands
}

func (p *Player) NonLandPermanents() []*Permanent {
	other := []*Permanent{}
	for _, perm := range p.Board {
		if !perm.IsLand && !perm.IsEnchantCreature {
			other = append(other, perm)
		}
	}
	return other
}

func (p *Player) AvatarString(position int, gameWidth int) string {
	playerString := ""
	for x := 0; x < (gameWidth-len(playerString))/2; x++ {
		playerString += " "
	}
	playerString += fmt.Sprintf("<Life: %d> Player %d <Mana: %d>", p.Life, position, p.ColorlessManaPool)
	return playerString
}

func PrintRowOfCards(cards []*Card, gameWidth int) {
	perms := []*Permanent{}
	for _, card := range cards {
		perms = append(perms, &Permanent{Card: card})
	}
	PrintRowOfPermanents(perms, gameWidth)
}

func PrintRowOfPermanents(perms []*Permanent, gameWidth int) {
	asciiImages := [][CARD_HEIGHT][CARD_WIDTH]string{}
	for _, perm := range perms {
		asciiImages = append(asciiImages, perm.AsciiImage(false))
	}
	for row := 0; row < CARD_HEIGHT; row++ {
		for x := 0; x < (gameWidth-len(perms)*(CARD_WIDTH+1))/2; x++ {
			fmt.Printf(" ")
		}
		for _, bitmap := range asciiImages {
			for _, char := range bitmap[row] {
				fmt.Printf(char)
			}
			fmt.Printf(" ")
		}
		fmt.Printf("%s", "\n")
	}
}

// GetCreature gets the first creature in play with the given name.
// It returns nil if there is no such creature in play.
func (p *Player) GetCreature(name CardName) *Permanent {
	for _, perm := range p.Board {
		if perm.Name == name {
			return perm
		}
	}
	return nil
}

func (p *Player) DealDamage(damage int) {
	p.Life -= damage
	p.DamageThisTurn += damage
}

func (p *Player) IsLegalTarget(c *Card, perm *Permanent) bool {
	if p != perm.Owner && perm.Hexproof {
		return false
	}
	for _, effect := range perm.Effects {
		if effect.Untargetable {
			return false
		}
		if p != perm.Owner && effect.Hexproof {
			return false
		}
	}
	return true
}

// HasLegalTarget returns whether the player has a legal target for casting this card.
func (p *Player) HasLegalTarget(c *Card) bool {
	for _, creature := range p.game.Creatures() {
		if p.IsLegalTarget(c, creature) {
			return true
		}
	}
	return false
}

func (p *Player) ResolveEffect(e *Effect) {
	if e.Summon != nil {
		perm := &Permanent{
			Card:  e.Summon,
			Owner: p,
		}
		perm.TurnPlayed = p.game.Turn
		p.Board = append(p.Board, perm)
		perm.HandleComingIntoPlay()
	}
}
