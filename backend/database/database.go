package database

// Database 数据库接口
type Database interface {
	Connect() error
	Close() error
	Query(sql string, args ...interface{}) (interface{}, error)
	Exec(sql string, args ...interface{}) (interface{}, error)
}

// User 用户表模型
type User struct {
	ID       string
	Username string
	Email    string
	Password string
	CreatedAt int64
	UpdatedAt int64
}

// Message 消息表模型
type Message struct {
	ID        string
	Sender    string
	Receiver  string
	Content   string
	Status    string
	CreatedAt int64
	UpdatedAt int64
}

// UserRepository 用户仓库
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
}

// MessageRepository 消息仓库
type MessageRepository interface {
	SaveMessage(msg *Message) error
	GetMessageByID(id string) (*Message, error)
	GetUserMessages(userID string, limit int) ([]*Message, error)
	UpdateMessageStatus(id string, status string) error
	DeleteMessage(id string) error
}
