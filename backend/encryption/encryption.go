package encryption

// Encryptor 加密接口
type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// Message 加密消息
type EncryptedMessage struct {
	ID        string
	Sender    string
	Receiver  string
	Content   string // 加密后的内容
	Timestamp int64
	IV        string // 初始化向量
}
