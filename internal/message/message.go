package message

const (
	OutMsg MessageType = iota
	ErrorMsg
)

type MessageType int

type Message struct {
	MessageType MessageType
	Body        string
}
