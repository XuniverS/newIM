package service

import (
	"im-system/server/internal/repository"
	"im-system/server/pkg/crypto"
)

// KeyService 密钥服务接口
type KeyService interface {
	GenerateKeys(userID int) (publicKey, privateKey string, err error)
	UploadPublicKey(userID int, publicKey string) error
	GetPublicKey(userID int) (string, error)
}

type keyService struct {
	repo     repository.KeyRepository
	userRepo repository.UserRepository
}

// NewKeyService 创建密钥服务实例
func NewKeyService(repo repository.KeyRepository, userRepo repository.UserRepository) KeyService {
	return &keyService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *keyService) GenerateKeys(userID int) (string, string, error) {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", "", err
	}

	// 检查是否已存在密钥
	exists, err := s.repo.Exists(userID)
	if err != nil {
		return "", "", err
	}
	if exists {
		return "", "", ErrKeyAlreadyExists
	}

	// 生成密钥对
	publicKey, privateKey, err := crypto.GenerateECCKeyPair()
	if err != nil {
		return "", "", err
	}

	// 保存公钥
	if err := s.repo.Save(userID, publicKey); err != nil {
		return "", "", err
	}

	return publicKey, privateKey, nil
}

func (s *keyService) UploadPublicKey(userID int, publicKey string) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.repo.Save(userID, publicKey)
}

func (s *keyService) GetPublicKey(userID int) (string, error) {
	return s.repo.Get(userID)
}

var ErrKeyAlreadyExists = &KeyError{"public key already exists"}

type KeyError struct {
	Message string
}

func (e *KeyError) Error() string {
	return e.Message
}
