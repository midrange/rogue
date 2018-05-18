/*
	A Cost currently accomodates colorless, Life for Phyrexian Spells,
	and Effects such as Quirion Ranger.

	TODO: expand to colored mana

*/

package game

import (
	"fmt"
)

type Cost struct {
	Colorless int
	Effect    *Effect
	Life      int // for Phyrexian
}

func (cc *Cost) String() string {
	if cc.Life > 0 {
		return fmt.Sprintf("%d (%d life)", cc.Colorless, cc.Life)
	} else {
		return fmt.Sprintf("%d", cc.Colorless)
	}
}
