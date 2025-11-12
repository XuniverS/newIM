package service

import (
	"errors"
	"time"

	"im-system/server/internal/config"
	"im-system/server/internal/model"
	"im-system/server/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserService 用户服务接口
type UserService interface {
	Register(username, password string) (string, int, error)
	Login(username, password string) (string, int, error)
	GetAllUsers(excludeUserID int) ([]model.UserPublicInfo, error)
	GetUserByID(userID int) (*model.User, error)
	GenerateToken(userID int, username string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type userService struct {
	repo   repository.UserRepository
	config *config.Config
}

// Claims JWT 声明
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewUserService 创建用户服务实例
func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo:   repo,
		config: cfg,
	}
}

func (s *userService) Register(username, password string) (string, int, error) {
	// 创建用户
	userID, err := s.repo.Create(username, password)
	if err != nil {
		return "", 0, errors.New("username already exists")
	}

	// 生成 token
	token, err := s.GenerateToken(userID, username)
	if err != nil {
		return "", 0, err
	}

	return token, userID, nil
}

func (s *userService) Login(username, password string) (string, int, error) {
	// 获取用户
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", 0, errors.New("invalid credentials")
	}

	// 验证密码
	if !s.repo.VerifyPassword(user.Password, password) {
		return "", 0, errors.New("invalid credentials")
	}

	// 生成 token
	token, err := s.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", 0, err
	}

	return token, user.ID, nil
}

func (s *userService) GetAllUsers(excludeUserID int) ([]model.UserPublicInfo, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []model.UserPublicInfo
	for _, user := range users {
		if user.ID != excludeUserID {
			result = append(result, model.UserPublicInfo{
				ID:       user.ID,
				Username: user.Username,
			})
		}
	}

	return result, nil
}

func (s *userService) GetUserByID(userID int) (*model.User, error) {
	return s.repo.GetByID(userID)
}

func (s *userService) GenerateToken(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *userService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
