package repository

import (
	"database/sql"
	"errors"

	"im-system/server/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(username, password string) (int, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(userID int) (*model.User, error)
	GetAll() ([]model.User, error)
	VerifyPassword(hashedPassword, password string) bool
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(username, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	err = r.db.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		username, string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(userID int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAll() ([]model.User, error) {
	rows, err := r.db.Query("SELECT id, username, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *userRepository) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
