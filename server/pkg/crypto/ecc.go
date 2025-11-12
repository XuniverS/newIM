package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/pem"
	"errors"
)

// GenerateECCKeyPair 生成 ECC 公私钥对 (P-256 曲线)
func GenerateECCKeyPair() (publicKeyPEM, privateKeyPEM string, err error) {
	// 使用 P-256 曲线生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// 编码私钥为 PEM 格式
	privateKeyBytes, err := encodePrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))

	// 编码公钥为 PEM 格式
	publicKeyBytes := encodePublicKey(&privateKey.PublicKey)
	publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	return publicKeyPEM, privateKeyPEM, nil
}

// EncryptWithPublicKey 使用 ECC 公钥加密（ECDH + AES-256-GCM）
func EncryptWithPublicKey(publicKeyPEM string, plaintext []byte) ([]byte, error) {
	// 解析接收者的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	recipientPublicKey, err := parsePublicKeyFromBytes(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 生成临时 ECDH 密钥对
	ephemeralPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// 执行 ECDH 密钥交换
	sharedSecret, err := ephemeralPrivateKey.ECDH(recipientPublicKey)
	if err != nil {
		return nil, err
	}

	// 使用 SHA-256 派生加密密钥
	hash := sha256.Sum256(sharedSecret)
	encryptionKey := hash[:]

	// 使用 AES-256-GCM 加密
	block256, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block256)
	if err != nil {
		return nil, err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// 返回：临时公钥 + 加密数据
	ephemeralPublicKeyBytes := ephemeralPrivateKey.PublicKey().Bytes()
	return append(ephemeralPublicKeyBytes, ciphertext...), nil
}

// DecryptWithPrivateKey 使用 ECC 私钥解密（ECDH + AES-256-GCM）
func DecryptWithPrivateKey(privateKeyPEM string, encryptedData []byte) ([]byte, error) {
	// 解析私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := parsePrivateKeyFromBytes(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 提取临时公钥（P-256 公钥为 65 字节）
	if len(encryptedData) < 65 {
		return nil, errors.New("invalid encrypted data format")
	}

	ephemeralPublicKeyBytes := encryptedData[:65]
	ciphertext := encryptedData[65:]

	// 解析临时公钥
	ephemeralPublicKey, err := ecdh.P256().NewPublicKey(ephemeralPublicKeyBytes)
	if err != nil {
		return nil, err
	}

	// 执行 ECDH 密钥交换
	sharedSecret, err := privateKey.ECDH(ephemeralPublicKey)
	if err != nil {
		return nil, err
	}

	// 使用 SHA-256 派生解密密钥
	hash := sha256.Sum256(sharedSecret)
	decryptionKey := hash[:]

	// 使用 AES-256-GCM 解密
	block256, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block256)
	if err != nil {
		return nil, err
	}

	// 提取 nonce 和密文
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// 辅助函数
func encodePrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	d := key.D.Bytes()
	if len(d) < 32 {
		d = append(make([]byte, 32-len(d)), d...)
	}
	return d, nil
}

func encodePublicKey(key *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(elliptic.P256(), key.X, key.Y)
}

func parsePublicKeyFromBytes(keyBytes []byte) (*ecdh.PublicKey, error) {
	return ecdh.P256().NewPublicKey(keyBytes)
}

func parsePrivateKeyFromBytes(keyBytes []byte) (*ecdh.PrivateKey, error) {
	if len(keyBytes) < 32 {
		keyBytes = append(make([]byte, 32-len(keyBytes)), keyBytes...)
	}
	return ecdh.P256().NewPrivateKey(keyBytes[:32])
}
