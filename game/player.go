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
	if card == NoCard {
		return
	}
	p.Hand = append(p.Hand, card)
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

// Returns possible actions when we can activate cards on the board.
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

		// TODO make actions unique, like don't allow two untapped Forests to both be cost targets
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
								Cost:   &Cost{Effect: costEffect},
								Owner:  p,
								Source: perm,
								Target: c,
								Type:   Activate,
							})
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
	answer := []*Action{}

	for _, name := range p.Hand {
		// Don't re-check playing duplicate cards
		if cardNames[name] {
			continue
		}
		cardNames[name] = true
		card := name.Card()

		if allowSorcerySpeed {
			if forHuman && card.IsEnchantment() && p.HasLegalTarget(card) {
				answer = append(answer, &Action{
					Type: ChooseTargetAndMana,
					Card: card,
				})
			}
			answer = p.appendActionsIfNonInstant(answer, card, forHuman)
		}

		if card.IsInstant() && (p.HasLegalTarget(card) || card.HasCreatureTargets() == false) {
			if forHuman {
				if p.CanPayCost(card.CastingCost) ||
					(card.Kicker != nil && p.CanPayCost(card.Kicker.Cost)) ||
					(card.PhyrexianCastingCost != nil && p.CanPayCost(card.PhyrexianCastingCost)) ||
					(card.AlternateCastingCost != nil && p.CanPayCost(card.AlternateCastingCost)) {
					answer = append(answer, &Action{
						Type: ChooseTargetAndMana,
						Card: card,
					})
				}
			} else {
				answer = p.appendActionsForInstant(answer, card)
			}
		}

	}

	return answer
}

// Appends actions to answer for an instant card.
func (p *Player) appendActionsForInstant(answer []*Action, card *Card) []*Action {
	if p.CanPayCost(card.CastingCost) {
		// TODO - add player targets - this assumes all instants target creatures for now
		for _, targetCreature := range p.game.Creatures() {
			if p.IsLegalTarget(card, targetCreature) {
				selectableLandCount := selectableLandCount(card)
				if selectableLandCount > 0 { // snap
					for i := 1; i <= len(p.game.Lands())-1; i++ {
						comb := combinations(makeRange(0, i), selectableLandCount)
						for _, c := range comb {
							selected := []*Permanent{}
							for _, index := range c {
								selected = append(selected, p.game.Lands()[index])
							}
							answer = append(answer, &Action{
								Type:     Play,
								Card:     card,
								Owner:    p,
								Selected: selected,
								Target:   targetCreature,
							})
						}
					}
				} else {
					answer = append(answer, &Action{
						Type:   Play,
						Card:   card,
						Owner:  p,
						Target: targetCreature,
					})
				}

			}
		}
	}
	if card.Kicker != nil && p.CanPayCost(card.Kicker.Cost) {
		for _, target := range p.game.Creatures() {
			if p.IsLegalTarget(card, target) {
				answer = append(answer, &Action{
					Type:       Play,
					Card:       card,
					Owner:      p,
					Target:     target,
					WithKicker: true,
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
					Owner:         p,
					Target:        target,
					WithPhyrexian: true,
				})
			}
		}
	}

	if card.AlternateCastingCost != nil && p.CanPayCost(card.AlternateCastingCost) {
		selectableLandCount := card.AlternateCastingCost.Effect.Selector.Count
		if selectableLandCount > 0 { // gush
			islands := p.landsOfSubtype(card.AlternateCastingCost.Effect.Selector.Subtype)
			for i := 1; i <= len(islands)-1; i++ {
				comb := combinations(makeRange(0, i), selectableLandCount)
				for _, c := range comb {
					selected := []*Permanent{}
					for _, index := range c {
						selected = append(selected, islands[index])
					}
					answer = append(answer, &Action{
						Type:          Play,
						Card:          card,
						Owner:         p,
						Selected:      selected,
						WithAlternate: true,
					})
				}
			}
		}

		answer = append(answer, &Action{Type: Play, Card: card, WithAlternate: true, Owner: p})
	}

	return answer
}

// Returns a list of lands of the given subtype.
func (p *Player) landsOfSubtype(subtype Subtype) []*Permanent {
	lands := []*Permanent{}
	for _, l := range p.Lands() {
		for _, st := range l.Subtype {
			if st == subtype {
				lands = append(lands, l)
				break
			}
		}
	}
	return lands
}

// Returns a count of lands the card effects would select.
func selectableLandCount(card *Card) int {
	count := 0
	for _, e := range card.Effects {
		if e.Selector != nil && e.Selector.Type == Land {
			count += e.Selector.Count
		}
	}
	return count
}

// Returns an array of ints from min to max.
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// Returns a list of lists of indexes of lands, based on https://play.golang.org/p/JEgfXR2zSH
func combinations(iterable []int, r int) [][]int {
	resultList := [][]int{}
	pool := iterable
	n := len(pool)

	if r > n {
		return resultList
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]int, r)
	for i, el := range indices {
		result[i] = pool[el]
	}

	resCopy := make([]int, r)
	copy(resCopy, result)
	resultList = append(resultList, resCopy)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			return resultList
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		newResCopy := make([]int, r)
		copy(newResCopy, result)
		resultList = append(resultList, newResCopy)
	}
	return resultList
}

// Appends actions to answer if card is a land, creature, or enchantment.
func (p *Player) appendActionsIfNonInstant(answer []*Action, card *Card, forHuman bool) []*Action {
	if card.IsLand() {
		if p.LandPlayedThisTurn == 0 {
			answer = append(answer, &Action{Type: Play, Card: card, Owner: p})
		}
	} else if !card.IsInstant() {
		if p.CanPayCost(card.CastingCost) {
			if card.IsCreature() {
				answer = append(answer, &Action{Type: Play, Card: card, Owner: p})
			} else if card.IsEnchantment() && p.HasLegalTarget(card) && !forHuman {
				for _, target := range p.game.Creatures() {
					answer = append(answer, &Action{
						Type:   Play,
						Card:   card,
						Owner:  p,
						Target: target,
					})
				}
			}
		}
		if card.PhyrexianCastingCost != nil && p.CanPayCost(card.PhyrexianCastingCost) {
			answer = append(answer, &Action{Type: Play, Card: card, WithPhyrexian: true, Owner: p})
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
						Target: attacker,
						With:   perm,
					})
				}
			}
		}
	}
	return answer
}

func (p *Player) RemoveCardForActionFromHand(action *Action) {
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
		p.game.Print()
		panic("XXX")
	}
	p.Hand = newHand
}

func (p *Player) PayCostsAndPutSpellOnStack(action *Action) {
	p.RemoveCardForActionFromHand(action)

	card := action.Card
	if !card.IsLand() {
		if action.WithKicker {
			p.PayCost(card.Kicker.Cost) // TODO use UpdatedEffectForAction when cardpool expands
		} else if action.WithAlternate {
			card.AlternateCastingCost.Effect = UpdatedEffectForAction(action, card.AlternateCastingCost.Effect)
			p.PayCost(card.AlternateCastingCost)
		} else if action.WithPhyrexian {
			p.PayCost(card.PhyrexianCastingCost) // TODO use UpdatedEffectForAction when cardpool expands
		} else {
			p.PayCost(card.CastingCost)
		}
	}
}

func (p *Player) PlayLand(action *Action) {
	p.RemoveCardForActionFromHand(action)
	card := action.Card
	p.game.newPermanent(card, p)
	p.LandPlayedThisTurn++
}

// Resolve an action to play a spell (non-land)
func (p *Player) ResolveSpell(action *Action) {
	card := action.Card

	for _, permanent := range p.Board {
		permanent.RespondToSpell()
	}

	if card.IsSpell() {
		p.CastSpell(card, action.Target, action)
		// TODO put spells (instants and sorceries) in graveyard (or exile)
	} else {
		// Non-spell (instant/sorcery) cards turn into permanents
		perm := p.game.newPermanent(card, p)

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
			p.ResolveEffect(UpdatedEffectForAction(a, e), nil)
			if target != nil {
				target.Plus1Plus1Counters += e.Plus1Plus1Counters // can be and often is 0 here
			}
		}
	}
	if c.Morbid != nil && (p.CreatureDied || p.Opponent().CreatureDied) && target != nil {
		target.Plus1Plus1Counters += c.Morbid.Plus1Plus1Counters
	}
}

func (p *Player) PayCostsAndPutAbilityOnStack(a *Action) {
	a.Source.PayForActivatedAbility(a.Cost, a.Target)
}

func (p *Player) ResolveActivatedAbility(a *Action) {
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
	} else if e.EffectType == ReturnToHand {
		// target is nil for rancor, or any effect of a permanent on itself
		effectedPermanent := e.Target
		if effectedPermanent == nil {
			effectedPermanent = perm
		}
		if effectedPermanent == nil {
			for _, selected := range e.Selected {
				p.RemoveFromBoard(selected)
				p.Hand = append(p.Hand, selected.Card.Name)
			}
		} else {
			p.RemoveFromBoard(effectedPermanent)
			p.Hand = append(p.Hand, effectedPermanent.Card.Name)
		}
	} else if e.EffectType == Untap {
		if e.Selector == nil { // nettle sentinel, or any effect of a permanent on itself
			perm.Tapped = false
		} else {
			for _, p := range e.Selected {
				p.Tapped = false
			}
		}
	} else if e.EffectType == AddMana {
		p.ColorlessManaPool += e.Colorless
	} else if e.EffectType == DrawCard {
		drawCount := 1
		if e.Selector != nil {
			drawCount = Max(drawCount, e.Selector.Count)
		}
		for i := 0; i < drawCount; i++ {
			p.Draw()
		}
	} else {
		panic("tried to resolve unknown effect")
	}
}

// Returns whether the player has the resources (life, mana, etc) to pay Cost.
func (p *Player) CanPayCost(c *Cost) bool {
	if c.Effect == nil {
		return p.Life >= c.Life && p.AvailableMana() >= c.Colorless
	} else {
		if c.Effect.EffectType == ReturnToHand {
			if c.Effect.Selector.Subtype != NoSubtype {
				count := Max(c.Effect.Selector.Count, 1)
				return len(p.landsOfSubtype(c.Effect.Selector.Subtype)) >= count
			}
		}
	}
	return false
}

// PayCost spends the resources for a Cost.
func (p *Player) PayCost(c *Cost) bool {

	// regular mana costs
	p.SpendMana(c.Colorless)

	// Phyrexian costs
	p.Life -= c.Life

	// costs like Gush
	if c.Effect != nil {
		p.ResolveEffect(c.Effect, nil)
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
		p.game.Print()
		panic("could not spend mana")
	}
}
