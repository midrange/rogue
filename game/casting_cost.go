/*
	A Casting Cost currently only accomodates colorless and Life for Phyrexian Spells.

	TODO: expand to colored mana

*/

package game

import (
	"fmt"
)

type CastingCost struct {
	Colorless int
	Life      int // for Phyrexian
}

func (cc *CastingCost) String() string {
	if cc.Life > 0 {
		return fmt.Sprintf("%d (%d life)", cc.Colorless, cc.Life)
	} else {
		return fmt.Sprintf("%d", cc.Colorless)
	}
}
