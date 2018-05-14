package game

import ()

type Modifier struct {
	Cost              int // when a Modifier is a kicker, it has a Cost
	Hexproof          bool
	Power             int
	PowerCounters     int
	Toughness         int
	ToughnessCounters int
	Untargetable      bool
}
