package types

type Move uint8

const (
	MoveUp Move = iota + 1
	MoveDown
	MoveLeft
	MoveRight
)
