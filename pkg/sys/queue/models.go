package queue

type Topic string
type MessageType string

type Message struct {
	Type MessageType
	Data []byte
}
