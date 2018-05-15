/*
	A Casting Cost currently only accomodates colorless and Life for Phyrexian Spells.

	TODO: expand to colored mana

*/

package game

import ()

type CastingCost struct {
	Colorless int
	Life      int // for Phyrexian
}
