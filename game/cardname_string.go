// Code generated by "stringer -type=CardName"; DO NOT EDIT.

package game

import "strconv"

const _CardName_name = "ForestGrizzlyBearsNettleSentinelRancorVinesOfVastwood"

var _CardName_index = [...]uint8{0, 6, 18, 32, 38, 53}

func (i CardName) String() string {
	if i < 0 || i >= CardName(len(_CardName_index)-1) {
		return "CardName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CardName_name[_CardName_index[i]:_CardName_index[i+1]]
}
