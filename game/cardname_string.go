// Code generated by "stringer -type=CardName"; DO NOT EDIT.

package game

import "strconv"

const _CardName_name = "NoCardBurningTreeEmissaryEldraziSpawnTokenForestGrizzlyBearsHungerOfTheHowlpackNestInvaderNettleSentinelQuirionRangerRancorSilhanaLedgewalkerSkarrganPitskulkVaultSkirgeVinesOfVastwood"

var _CardName_index = [...]uint8{0, 6, 25, 42, 48, 60, 79, 90, 104, 117, 123, 141, 157, 168, 183}

func (i CardName) String() string {
	if i < 0 || i >= CardName(len(_CardName_index)-1) {
		return "CardName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CardName_name[_CardName_index[i]:_CardName_index[i+1]]
}
