// Code generated by "stringer -type=EffectType"; DO NOT EDIT.

package game

import "strconv"

const _EffectType_name = "AddManaDrawCardReturnToHandUntap"

var _EffectType_index = [...]uint8{0, 7, 15, 27, 32}

func (i EffectType) String() string {
	if i < 0 || i >= EffectType(len(_EffectType_index)-1) {
		return "EffectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EffectType_name[_EffectType_index[i]:_EffectType_index[i+1]]
}
