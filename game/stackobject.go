package game

import ()

type StackObject struct {
	Type     ActionType
	Card     *Card // for spell-based stack objects
	Cost     *Cost
	Kicker   *Effect
	Player   *Player
	Selected []*Permanent
	Source   *Permanent
	Target   *Permanent // a target that is a Permanent (players not yet handled)
}
