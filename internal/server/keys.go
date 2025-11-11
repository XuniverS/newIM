package server

import (
	"net/http"
	"strconv"

	"im-system/internal/crypto"
	"im-system/internal/db"

	"github.com/gin-gonic/gin"
)

type GenerateKeysResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type UploadPublicKeyRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
}

func (s *Server) handleGenerateKeys(c *gin.Context) {
	userID := s.getUserIDFromContext(c)

	// 检查是否已经存在公钥
	exists, err := db.PublicKeyExists(s.db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Public key already exists"})
		return
	}

	// 生成 RSA 密钥对
	publicKey, privateKey, err := crypto.GenerateRSAKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate keys"})
		return
	}

	// 保存公钥到数据库
	if err := db.SavePublicKey(s.db, userID, publicKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save public key"})
		return
	}

	c.JSON(http.StatusOK, GenerateKeysResponse{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
}

func (s *Server) handleUploadPublicKey(c *gin.Context) {
	userID := s.getUserIDFromContext(c)

	var req UploadPublicKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 保存公钥
	if err := db.SavePublicKey(s.db, userID, req.PublicKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save public key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Public key uploaded successfully"})
}

func (s *Server) handleGetPublicKey(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	publicKey, err := db.GetPublicKey(s.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Public key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"public_key": publicKey})
}
