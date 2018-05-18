/*
	A TargetType is what a spell or effect can target.

	This might be something like "creature you control" or "a player."
*/

package game

import (
	"fmt"
)

//go:generate stringer -type=PlayerTargetType
type PlayerTargetType int

const (
	SamePlayer PlayerTargetType = iota
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
const (
	Artifact Type = iota
	Creature
	Enchantment
	Instant
	Land
	Planeswalker
	Sorcery
	Tribal
)

type TargetType struct {
	ControlledBy PlayerTargetType
	Supertype    Supertype
	Subtype      Subtype
	Type         Type
}

func (tt *TargetType) String() string {
	return fmt.Sprintf("%s, %s, %s  - controlled by %s", tt.Type, tt.Subtype, tt.Supertype, tt.ControlledBy)
}
