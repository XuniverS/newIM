package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"im-system/client/internal/config"
	"im-system/client/internal/model"
)

// ServerService 服务端通信服务接口
type ServerService interface {
	Register(username, password string) (*model.AuthResponse, error)
	Login(username, password string) (*model.AuthResponse, error)
	GetAllUsers(token string) ([]model.User, error)
	GetOnlineUsers(token string) ([]int, error)
	GetPublicKey(token string, userID int) (string, error)
	GenerateKeys(token string) (*model.KeyPair, error)
	SendMessage(token string, receiverID int, encryptedContent string) (int, error)
	GetUnreadMessages(token string) ([]model.Message, error)
	GetServerWSURL() string
}

type serverService struct {
	config *config.Config
	client *http.Client
}

// NewServerService 创建服务端通信服务实例
func NewServerService(cfg *config.Config) ServerService {
	return &serverService{
		config: cfg,
		client: &http.Client{},
	}
}

func (s *serverService) Register(username, password string) (*model.AuthResponse, error) {
	reqBody := model.AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := s.post("/api/auth/register", "", reqBody)
	if err != nil {
		return nil, err
	}

	var authResp model.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (s *serverService) Login(username, password string) (*model.AuthResponse, error) {
	reqBody := model.AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := s.post("/api/auth/login", "", reqBody)
	if err != nil {
		return nil, err
	}

	var authResp model.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (s *serverService) GetAllUsers(token string) ([]model.User, error) {
	resp, err := s.get("/api/users", token)
	if err != nil {
		return nil, err
	}

	var result struct {
		Users []model.User `json:"users"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.Users, nil
}

func (s *serverService) GetOnlineUsers(token string) ([]int, error) {
	resp, err := s.get("/api/users/online", token)
	if err != nil {
		return nil, err
	}

	var result struct {
		OnlineUsers []int `json:"online_users"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.OnlineUsers, nil
}

func (s *serverService) GetPublicKey(token string, userID int) (string, error) {
	resp, err := s.get(fmt.Sprintf("/api/keys/%d", userID), token)
	if err != nil {
		return "", err
	}

	var result struct {
		PublicKey string `json:"public_key"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result.PublicKey, nil
}

func (s *serverService) GenerateKeys(token string) (*model.KeyPair, error) {
	resp, err := s.post("/api/keys/generate", token, nil)
	if err != nil {
		return nil, err
	}

	var keyPair model.KeyPair
	if err := json.Unmarshal(resp, &keyPair); err != nil {
		return nil, err
	}

	return &keyPair, nil
}

func (s *serverService) SendMessage(token string, receiverID int, encryptedContent string) (int, error) {
	reqBody := map[string]interface{}{
		"receiver_id": receiverID,
		"content":     encryptedContent,
	}

	resp, err := s.post("/api/messages/send", token, reqBody)
	if err != nil {
		return 0, err
	}

	var result struct {
		MessageID int `json:"message_id"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, err
	}

	return result.MessageID, nil
}

func (s *serverService) GetUnreadMessages(token string) ([]model.Message, error) {
	resp, err := s.get("/api/messages/unread", token)
	if err != nil {
		return nil, err
	}

	var result struct {
		Messages []model.Message `json:"messages"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result.Messages, nil
}

func (s *serverService) GetServerWSURL() string {
	return s.config.GetServerWSURL() + "/api/ws"
}

// 辅助方法
func (s *serverService) get(path, token string) ([]byte, error) {
	url := s.config.GetServerURL() + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	return body, nil
}

func (s *serverService) post(path, token string, data interface{}) ([]byte, error) {
	url := s.config.GetServerURL() + path

	var reqBody []byte
	var err error
	if data != nil {
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	return body, nil
}
