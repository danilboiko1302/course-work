package types

type Action uint16

const (
	ErrorAction Action = iota
	Login
	Logout
	GetUsers
	SetAdmin
	GetField
	Up
	Down
	Right
	Left
	SendMessage
	GetHistory
)

type ActionFront uint16

const (
	ErrorActionFront ActionFront = iota
	NewUser
	UserLeft
	AllUsers
	LogedIn
	Field
	Lost
	NewAdmin
	GetMessage
	GetMessages
)

type MessageFront struct {
	Action ActionFront `json:"action"`
	Data   string      `json:"data"`
}

type Message struct {
	Action Action `json:"action"`
	Data   string `json:"data"`
}
