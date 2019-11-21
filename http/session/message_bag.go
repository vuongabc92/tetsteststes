package session

type MessageBag struct {
	Messages map[string]string
}

func (m MessageBag) Add(key string, msg string) {
	m.Messages[key] = msg
}

func (m MessageBag) Has(key string) bool {
	if m.IsError() {
		_, ok := m.Messages[key]
		return ok
	}

	return false
}

func (m MessageBag) Get(key string) string {
	if m.Has(key) {
		return m.Messages[key]
	}

	return ""
}

func (m MessageBag) GetMessages() map[string]string {
	return m.Messages
}

func (m MessageBag) IsError() bool {
	if len(m.Messages) > 0 {
		return true
	}

	return false
}

func NewMessageBag() *MessageBag {
	return &MessageBag{Messages: make(map[string]string)}
}
