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
	Hand               []CardName
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
		Hand:  []CardName{},
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

func (p *Player) AddToHand(c CardName) {
	if c == NoCard {
		return
	}
	p.Hand = append(p.Hand, c)
}

func (p *Player) AvailableMana() int {
	answer := 0
	for _, card := range p.Board {
		if card.IsLand() && !card.Tapped {
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
	for _, perm := range p.Board {
		perm.Damage = 0
		perm.TemporaryEffects = []*Effect{}
		perm.ActivatedThisTurn = false
	}
	p.LandPlayedThisTurn = 0
	p.DamageThisTurn = 0
	p.CreatureDied = false
	p.EndPhase()
}

func (p *Player) Creatures() []*Permanent {
	answer := []*Permanent{}
	for _, card := range p.Board {
		if card.IsCreature() {
			answer = append(answer, card)
		}
	}
	return answer
}

func (p *Player) SendToGraveyard(perm *Permanent) {
	p.RemoveFromBoard(perm)
	if perm.EntersGraveyardEffect != nil {
		p.ResolveEffect(perm.EntersGraveyardEffect, perm)
	}

	if perm.IsCreature() {
		p.CreatureDied = true
	}

	perm.TemporaryEffects = []*Effect{}
	for _, aura := range perm.Auras {
		if aura.EnchantedPermanentDiesEffect != nil {
			p.ResolveEffect(aura.EnchantedPermanentDiesEffect, aura)
		}
		p.SendToGraveyard(aura)
	}
}

func (p *Player) RemoveFromBoard(perm *Permanent) {
	newBoard := []*Permanent{}
	for _, permanent := range p.Board {
		if permanent != perm {
			newBoard = append(newBoard, permanent)
		}
	}
	p.Board = newBoard
}

// Returns possible actions when we can activate cards on he board.
func (p *Player) ActivatedAbilityActions(allowSorcerySpeed bool, forHuman bool) []*Action {
	permNames := make(map[CardName]bool)
	answer := []*Action{}

	for _, perm := range p.Board { // TODO could be opponent's board for some actions, e.g. Warmonger
		// Don't re-check playing duplicate actions
		if permNames[perm.Name] {
			continue
		}
		if perm.ActivatedThisTurn {
			continue
		}
		permNames[perm.Name] = true

		// TODO make actions unique, like don't allow two untaped Forests to both be cost targets
		if perm.ActivatedAbility != nil {
			effect := perm.ActivatedAbility
			landsForCost := []*Permanent{}
			if effect.Cost.Effect != nil && effect.Cost.Effect.Selector.Subtype != NoSubtype {
				for _, l := range p.Lands() {
					if l.HasSubtype(effect.Cost.Effect.Selector.Subtype) {
						landsForCost = append(landsForCost, l)
					}
				}
			}

			if effect.Selector.Type == Creature { // TODO lands etc
				for _, c := range p.Creatures() {
					for _, land := range landsForCost {
						costEffect := effect.Cost.Effect
						costEffect.SelectedForCost = land
						answer = append(answer,
							&Action{
								Type:   Activate,
								Source: perm,
								Cost:   &Cost{Effect: costEffect},
								Target: c})
					}
				}
			}
		}
	}
	return answer
}

// Returns possible actions when we can play a card from hand, including passing.
func (p *Player) PlayActions(allowSorcerySpeed bool, forHuman bool) []*Action {
	cardNames := make(map[CardName]bool)
	answer := []*Action{&Action{Type: Pass}}

	for _, name := range p.Hand {
		// Don't re-check playing duplicate cards
		if cardNames[name] {
			continue
		}
		cardNames[name] = true
		card := name.Card()

		if allowSorcerySpeed {
			if card.IsLand() {
				if p.LandPlayedThisTurn == 0 {
					answer = append(answer, &Action{Type: Play, Card: card})
				}
			} else if !card.IsInstant() {
				if p.CanPayCost(card.CastingCost) {
					if card.IsCreature() {
						answer = append(answer, &Action{Type: Play, Card: card})
					}
					if card.IsEnchantment() && p.CanPayCost(card.CastingCost) && p.HasLegalTarget(card) {
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
				}
				if card.PhyrexianCastingCost != nil && p.CanPayCost(card.PhyrexianCastingCost) {
					answer = append(answer, &Action{Type: Play, Card: card, WithPhyrexian: true})
				}
			}
		}

		if card.IsInstant() && p.HasLegalTarget(card) {
			if forHuman {
				if p.CanPayCost(card.CastingCost) ||
					(card.Kicker != nil && p.CanPayCost(card.Kicker.Cost)) ||
					(card.PhyrexianCastingCost != nil) && p.CanPayCost(card.PhyrexianCastingCost) {
					answer = append(answer, &Action{
						Type: ChooseTargetAndMana,
						Card: card,
					})
				}
			} else { // TODO - add player targets - this assumes all instants target creatures for now
				if p.CanPayCost(card.CastingCost) {
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
				if card.Kicker != nil && p.CanPayCost(card.Kicker.Cost) {
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

				if card.PhyrexianCastingCost != nil && p.CanPayCost(card.PhyrexianCastingCost) {
					for _, target := range p.game.Creatures() {
						if p.IsLegalTarget(card, target) {
							answer = append(answer, &Action{
								Type:          Play,
								Card:          card,
								Target:        target,
								WithPhyrexian: true,
							})
						}
					}
				}
			}
		} else if card.IsInstant() {
			if card.AlternateCastingCost != nil && p.CanPayCost(card.AlternateCastingCost) {
				answer = append(answer, &Action{Type: Play, Card: card, WithAlternate: true})
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
		if card.IsCreature() && !card.Attacking && !card.Tapped && card.TurnPlayed != p.game.Turn {
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
		if perm.Blocking == nil && !perm.Tapped && perm.IsCreature() {
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
	newHand := []CardName{}
	found := false
	for _, c := range p.Hand {
		if !found && c == card.Name {
			found = true
			continue
		}
		newHand = append(newHand, c)
	}
	if !found {
		log.Printf("could not play card %+v from hand %+v", card, p.Hand)
		panic("XXX")
	}
	p.Hand = newHand

	if !card.IsLand() {
		if action.WithKicker {
			p.PayCost(card.Kicker.Cost)
		} else if action.WithAlternate {
			p.PayCost(card.AlternateCastingCost)
		} else if action.WithPhyrexian {
			p.PayCost(card.PhyrexianCastingCost)
		} else {
			p.PayCost(card.CastingCost)
		}
		for _, permanent := range p.Board {
			permanent.RespondToSpell()
		}
	}

	if card.IsSpell() {
		p.CastSpell(card, action.Target, action)
		// TODO put spells (instants and sorceries) in graveyard (or exile)
	} else {
		// Non-spell (instant/sorcery) cards turn into permanents
		perm := p.game.newPermanent(card, p)

		if card.IsLand() {
			p.LandPlayedThisTurn++
		}

		if card.IsEnchantCreature() {
			action.Target.Auras = append(action.Target.Auras, perm)
		}
	}

}

func (p *Player) CastSpell(c *Card, target *Permanent, a *Action) {
	if c.AddsTemporaryEffect {
		for _, e := range c.Effects {
			target.TemporaryEffects = append(target.TemporaryEffects, UpdatedEffectForAction(a, e))
		}
	} else if c.Effects != nil {
		for _, e := range c.Effects {
			if target == nil {
				p.ResolveEffect(UpdatedEffectForAction(a, e), nil)
			} else {
				target.Plus1Plus1Counters += e.Plus1Plus1Counters // this is usually 0
				p.ResolveEffect(UpdatedEffectForAction(a, e), nil)
			}
		}
	}
	if c.Morbid != nil && (p.CreatureDied || p.Opponent().CreatureDied) && target != nil {
		target.Plus1Plus1Counters += c.Morbid.Plus1Plus1Counters
	}
}

func (p *Player) ActivateAbility(a *Action) {
	a.Source.ActivateAbility(a.Cost, a.Target)
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
		if perm.IsLand() {
			lands = append(lands, perm)
		}
	}
	return lands
}

func (p *Player) NonLandPermanents() []*Permanent {
	other := []*Permanent{}
	for _, perm := range p.Board {
		if !perm.IsLand() && !perm.IsEnchantment() {
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

func PrintRowOfCards(cards []CardName, gameWidth int) {
	perms := []*Permanent{}
	for _, name := range cards {
		perms = append(perms, &Permanent{Card: name.Card()})
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
	for _, effect := range perm.TemporaryEffects {
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

func (p *Player) ResolveEffect(e *Effect, perm *Permanent) {
	if e.Condition != nil {
		if e.Condition.ControlAnother != NoCard {
			controlsOne := false
			for _, boardPerm := range p.Board {
				if boardPerm.Name == e.Condition.ControlAnother && boardPerm.Id != perm.Id {
					controlsOne = true
					break
				}
			}
			if !controlsOne {
				return
			}
		} else {
			panic("unhandled Condition in ResolveEffect")
		}
	}
	if e.Summon != NoCard {
		p.game.newPermanent(e.Summon.Card(), p)
		return
	} else if e.EffectType == ReturnToHand {
		fmt.Println("effect is ", e)
		if e.Target == nil { // quirion ranger, gush?
			p.RemoveFromBoard(perm)
			p.Hand = append(p.Hand, perm.Card.Name)
		} else {
			p.RemoveFromBoard(e.Target)
			p.Hand = append(p.Hand, e.Target.Card.Name)
		}
		return
	} else if e.EffectType == Untap {
		if e.Selector == nil {
			perm.Tapped = false
		} else {
			if e.Selector.Subtype != NoSubtype {
				count := Max(e.Selector.Count, 1)
				for _, l := range p.Lands() {
					for _, st := range l.Subtype {
						if st == e.Selector.Subtype {
							l.Tapped = false
							count--
							break
						}
					}
					if count == 0 {
						break
					}
				}
			} else if e.Selector.Type == Creature {
				count := Max(e.Selector.Count, 1)
				for _, c := range p.game.Creatures() {
					c.Tapped = false
					if count == 0 {
						break
					}
				}
			}
		}
		return
	} else if e.EffectType == AddMana {
		p.ColorlessManaPool += e.Colorless
	} else if e.EffectType == DrawCard {
		for i := 0; i < Max(1, e.EffectCount); i++ {
			p.Draw()
		}
	} else {
		panic("tried to resolve unknown effect")
	}
}

// CanPayCost returns whether the player has the resources (life, mana, etc) to pay Cost.
func (p *Player) CanPayCost(c *Cost) bool {
	if c.Effect == nil {
		return p.Life >= c.Life && p.AvailableMana() >= c.Colorless
	} else {
		if c.Effect.EffectType == ReturnToHand {
			if c.Effect.Selector.Subtype != NoSubtype {
				count := Max(c.Effect.Selector.Count, 1)
				for _, l := range p.Lands() {
					for _, st := range l.Subtype {
						if st == c.Effect.Selector.Subtype {
							count--
							break
						}
					}
					if count == 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

// PayCost spends the resources for a Cost
func (p *Player) PayCost(c *Cost) bool {

	// regular mana costs
	p.SpendMana(c.Colorless)

	// Phyrexian costs
	p.Life -= c.Life

	// costs like Gush
	if c.Effect != nil {
		if c.Effect.EffectType == ReturnToHand {
			if c.Effect.Selector.Subtype != NoSubtype {
				count := Max(c.Effect.Selector.Count, 1)
				for _, l := range p.Lands() {
					for _, st := range l.Subtype {
						if st == c.Effect.Selector.Subtype {
							p.RemoveFromBoard(l)
							p.Hand = append(p.Hand, l.Card.Name)
							count--
							break
						}
					}
					if count == 0 {
						break
					}
				}
			}
		}
	}
	return false
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
		if card.IsLand() && !card.Tapped {
			card.Tapped = true
			amount -= 1
		}
	}
	if amount > 0 {
		panic("could not spend mana")
	}
}
