/*
	A Modifier is a struct with values that affect a Card's or Player's state.

	A Modifier can be the property of a Card. For example, when the Card is a Giant Growth,
	its modifier property would be &Modifier{power:3, toughness:3}.

	A Modifer can also be the Morbid or Kicker property of a Card, to designate a
	Modifier that only happens under special circumstances.

*/

package game

import ()

type Modifier struct {
	CastingCost       *CastingCost // when a Modifier is a kicker, it has a Cost
	Hexproof          bool
	Power             int
	PowerCounters     int
	Toughness         int
	ToughnessCounters int
	Untargetable      bool
}
