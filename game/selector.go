/*
	A Selector is what a spell or effect can target or use as a cost.

	This might be something like "creature you control" or "a player."
*/

package game

import (
	"fmt"
)

//go:generate stringer -type=PlayerSelector
type PlayerSelector int

const (
	SamePlayer PlayerSelector = iota
	OpposingPlayer
)

// https://mtg.gamepedia.com/Subtype
//go:generate stringer -type=Supertype
type Supertype int

const (
	Basic Supertype = iota
	Legendary
	Snow
	World
)

//go:generate stringer -type=Subtype
type Subtype int

// For example, Cat, Goblin, Bird, and Elf are creature subtypes.
// Ajani and Jace are planeswalker subtypes.
// Equipment, Aura, Trap, Arcane, are more subtypes.
const (
	NoSubtype Subtype = iota
	LandForest
	LandIsland
	LandMountain
	LandPlains
	LandSwamp
)

//go:generate stringer -type=Type
type Type int

// Some card types appear only on cards used in variants such as Planechase and Archenemy.
// Phenomenon, Vanguards, Schemes
// the Type Spell denotes all other Types except Land
const (
	Artifact Type = iota
	Creature
	Enchantment
	Instant
	Land
	Planeswalker
	Sorcery
	Tribal
	Spell
)

type Selector struct {
	Count        int
	ControlledBy PlayerSelector
	Supertype    Supertype
	Subtype      Subtype
	Targetted    bool
	Type         Type
}

func (s *Selector) String() string {
	return fmt.Sprintf("%s, %s, %s  - controlled by %s", s.Type, s.Subtype, s.Supertype, s.ControlledBy)
}
