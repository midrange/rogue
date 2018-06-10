package game

import ()

type LandState struct {
	ForestTappedCount   int
	ForestUntappedCount int
	IslandTappedCount   int
	IslandUntappedCount int
	TargetId            PermanentId
	TargetTapped        bool
}
