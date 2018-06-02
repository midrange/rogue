/*
	A Condition currently accommodates Faerie Miscreant and Spellstutter Sprite.
*/

package game

import (
	"fmt"
)

type Condition struct {
	ControlAnother       CardName
	ConvertedManaCostLTE Subtype
}

func (c *Condition) String() string {
	return fmt.Sprintf("Condition: controls %s", c.ControlAnother)
}
