package service

import (
	"encoding/base64"

	"im-system/client/pkg/crypto"
)

// CryptoService 加密服务接口
type CryptoService interface {
	GenerateKeyPair() (publicKey, privateKey string, err error)
	Encrypt(publicKeyPEM string, plaintext string) (string, error)
	Decrypt(privateKeyPEM string, ciphertext string) (string, error)
}

type cryptoService struct{}

// NewCryptoService 创建加密服务实例
func NewCryptoService() CryptoService {
	return &cryptoService{}
}

func (s *cryptoService) GenerateKeyPair() (string, string, error) {
	return crypto.GenerateECCKeyPair()
}

func (s *cryptoService) Encrypt(publicKeyPEM string, plaintext string) (string, error) {
	encrypted, err := crypto.EncryptWithPublicKey(publicKeyPEM, []byte(plaintext))
	if err != nil {
		return "", err
	}

	// 返回Base64编码的密文
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (s *cryptoService) Decrypt(privateKeyPEM string, ciphertext string) (string, error) {
	// 解码Base64
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	decrypted, err := crypto.DecryptWithPrivateKey(privateKeyPEM, encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
