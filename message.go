package wire

type Message struct {
	Id   string
	Body interface{}
}

func NewMessage(id string, body interface{}) *Message {
	return &Message{id, body}
}
