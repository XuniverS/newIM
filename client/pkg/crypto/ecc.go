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

// GenerateECCKeyPair 生成 ECC 公私钥对
func GenerateECCKeyPair() (publicKeyPEM, privateKeyPEM string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes, err := encodePrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))

	publicKeyBytes := encodePublicKey(&privateKey.PublicKey)
	publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	return publicKeyPEM, privateKeyPEM, nil
}

// EncryptWithPublicKey 使用公钥加密
func EncryptWithPublicKey(publicKeyPEM string, plaintext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	recipientPublicKey, err := parsePublicKeyFromBytes(block.Bytes)
	if err != nil {
		return nil, err
	}

	ephemeralPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := ephemeralPrivateKey.ECDH(recipientPublicKey)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(sharedSecret)
	encryptionKey := hash[:]

	block256, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block256)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	ephemeralPublicKeyBytes := ephemeralPrivateKey.PublicKey().Bytes()
	return append(ephemeralPublicKeyBytes, ciphertext...), nil
}

// DecryptWithPrivateKey 使用私钥解密
func DecryptWithPrivateKey(privateKeyPEM string, encryptedData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	privateKey, err := parsePrivateKeyFromBytes(block.Bytes)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < 65 {
		return nil, errors.New("invalid encrypted data format")
	}

	ephemeralPublicKeyBytes := encryptedData[:65]
	ciphertext := encryptedData[65:]

	ephemeralPublicKey, err := ecdh.P256().NewPublicKey(ephemeralPublicKeyBytes)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := privateKey.ECDH(ephemeralPublicKey)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(sharedSecret)
	decryptionKey := hash[:]

	block256, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block256)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

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
