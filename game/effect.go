package game

import ()

type Effect struct {
	BasePower     int
	BaseToughness int
	Card          *Card
	Untargetable  bool
}
