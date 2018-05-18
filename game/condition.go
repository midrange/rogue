/*
	A Condition currently accommodates just InPlay for a given CardName,
	so we can do Faerie Miscreant.
*/

package game

import (
	"fmt"
)

type Condition struct {
	ControlAnother CardName
}

func (c *Condition) String() string {
	return fmt.Sprintf("Condition: controls %s", c.ControlAnother)
}
