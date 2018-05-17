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
//go:generate stringer -type=MTGSubtype
type MTGSubtype int

const (
	BasicForest MTGSubtype = iota
	BasicIsland
)

//go:generate stringer -type=MTGType
type MTGType int

const (
	Creature MTGType = iota
	Land
)

type TargetType struct {
	ControlledBy PlayerTargetType
	Subtype      MTGSubtype
	Type         MTGType
}

func (tt *TargetType) String() string {
	return fmt.Sprintf("%s, %s - controlled by %s", tt.Type, tt.Subtype, tt.ControlledBy)
}
