package game

import ()

type Effect struct {
	Card         *Card
	Hexproof     bool
	Power        int
	Toughness    int
	Untargetable bool
}
