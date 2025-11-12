package service

import (
	"im-system/server/internal/model"
	"im-system/server/internal/repository"
)

// MessageService 消息服务接口
type MessageService interface {
	SendMessage(senderID, receiverID int, encryptedContent string) (int, error)
	GetUnreadMessages(userID int) ([]model.MessageDTO, error)
	MarkAsRead(messageID int) error
	GetConversation(userID1, userID2 int, limit int) ([]model.MessageDTO, error)
}

type messageService struct {
	repo     repository.MessageRepository
	userRepo repository.UserRepository
}

// NewMessageService 创建消息服务实例
func NewMessageService(repo repository.MessageRepository, userRepo repository.UserRepository) MessageService {
	return &messageService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *messageService) SendMessage(senderID, receiverID int, encryptedContent string) (int, error) {
	// 验证接收者存在
	_, err := s.userRepo.GetByID(receiverID)
	if err != nil {
		return 0, err
	}

	// 保存消息
	return s.repo.Save(senderID, receiverID, encryptedContent)
}

func (s *messageService) GetUnreadMessages(userID int) ([]model.MessageDTO, error) {
	messages, err := s.repo.GetUnread(userID)
	if err != nil {
		return nil, err
	}

	var result []model.MessageDTO
	for _, msg := range messages {
		result = append(result, model.MessageDTO{
			ID:               msg.ID,
			SenderID:         msg.SenderID,
			ReceiverID:       msg.ReceiverID,
			EncryptedContent: msg.EncryptedContent,
			IsRead:           msg.IsRead,
			CreatedAt:        msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}

func (s *messageService) MarkAsRead(messageID int) error {
	return s.repo.MarkAsRead(messageID)
}

func (s *messageService) GetConversation(userID1, userID2 int, limit int) ([]model.MessageDTO, error) {
	messages, err := s.repo.GetConversation(userID1, userID2, limit)
	if err != nil {
		return nil, err
	}

	var result []model.MessageDTO
	for _, msg := range messages {
		result = append(result, model.MessageDTO{
			ID:               msg.ID,
			SenderID:         msg.SenderID,
			ReceiverID:       msg.ReceiverID,
			EncryptedContent: msg.EncryptedContent,
			IsRead:           msg.IsRead,
			CreatedAt:        msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}
