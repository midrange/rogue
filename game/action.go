package game

import ()

type Action interface {
	Type() ActionType
}

type ActionType int

const (
	Pass ActionType = iota
	PlayLand
	Cast
	Attack
)
