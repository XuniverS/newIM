package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

// GenerateRSAKeyPair 生成 RSA 密钥对
// 返回：私钥PEM格式字符串、公钥PEM格式字符串、错误
func GenerateRSAKeyPair(bits int) (string, string, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("生成私钥失败: %v", err)
	}

	// 将私钥转换为 PEM 格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 将公钥转换为 PEM 格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("序列化公钥失败: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// RSAEncrypt 使用公钥加密字符串
// 参数：明文字符串、公钥PEM格式字符串
// 返回：Base64编码的密文、错误
func RSAEncrypt(plaintext string, publicKeyPEM string) (string, error) {
	// 解析 PEM 格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("解析公钥失败")
	}

	// 解析公钥
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析公钥失败: %v", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("不是有效的 RSA 公钥")
	}

	// 使用 OAEP 填充方式加密
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(plaintext),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}

	// 将密文转换为 Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSADecrypt 使用私钥解密字符串
// 参数：Base64编码的密文、私钥PEM格式字符串
// 返回：明文字符串、错误
func RSADecrypt(ciphertext string, privateKeyPEM string) (string, error) {
	// 解码 Base64 密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %v", err)
	}

	// 解析 PEM 格式的私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("解析私钥失败")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %v", err)
	}

	// 使用 OAEP 填充方式解密
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		ciphertextBytes,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	return string(plaintext), nil
}

func main() {
	fmt.Println("=== RSA 非对称加密演示 ===\n")

	// 1. 生成密钥对（2048位）
	fmt.Println("1. 生成 RSA 密钥对...")
	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatalf("生成密钥对失败: %v", err)
	}
	fmt.Println("✓ 密钥对生成成功\n")

	// 打印公钥和私钥（实际应用中私钥需要保密）
	fmt.Println("公钥:")
	fmt.Println(publicKey)
	fmt.Println("私钥:")
	fmt.Println(privateKey)

	// 2. 加密消息
	originalMessage := "Hello, 这是一条需要加密的消息！"
	fmt.Printf("2. 原始消息: %s\n", originalMessage)

	encryptedMessage, err := RSAEncrypt(originalMessage, publicKey)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}
	fmt.Printf("✓ 加密后的消息 (Base64): %s\n\n", encryptedMessage)

	// 3. 解密消息
	fmt.Println("3. 解密消息...")
	decryptedMessage, err := RSADecrypt(encryptedMessage, privateKey)
	if err != nil {
		log.Fatalf("解密失败: %v", err)
	}
	fmt.Printf("✓ 解密后的消息: %s\n\n", decryptedMessage)

	// 4. 验证结果
	if originalMessage == decryptedMessage {
		fmt.Println("✅ 加密解密成功！原始消息和解密后的消息一致")
	} else {
		fmt.Println("❌ 加密解密失败！消息不一致")
	}
}