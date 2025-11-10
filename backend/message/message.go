package message

// Message 消息结构
type Message struct {
	ID        string
	Sender    string
	Receiver  string
	Content   string
	Timestamp int64
	Status    string // pending, sent, delivered, read
	IsEncrypted bool
}

// MessageStore 消息存储接口
type MessageStore interface {
	SaveMessage(msg *Message) error
	GetMessage(id string) (*Message, error)
	GetUserMessages(userID string, limit int) ([]*Message, error)
	UpdateMessageStatus(id string, status string) error
	DeleteMessage(id string) error
}

// MessageService 消息服务
type MessageService struct {
	store MessageStore
}

func NewMessageService(store MessageStore) *MessageService {
	return &MessageService{store: store}
}

func (ms *MessageService) SendMessage(msg *Message) error {
	return ms.store.SaveMessage(msg)
}
