// Code generated by "stringer -type=ActionType"; DO NOT EDIT.

package game

import "strconv"

const _ActionType_name = "PassPlayDeclareAttackAttackBlockUseForManaChooseTargetAndManaActivateOfferToResolveNextOnStackResolveNextOnStack"

var _ActionType_index = [...]uint8{0, 4, 8, 21, 27, 32, 42, 61, 69, 94, 112}

func (i ActionType) String() string {
	if i < 0 || i >= ActionType(len(_ActionType_index)-1) {
		return "ActionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActionType_name[_ActionType_index[i]:_ActionType_index[i+1]]
}
