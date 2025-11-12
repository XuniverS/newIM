package repository

import (
	"database/sql"
	"errors"
)

// KeyRepository 密钥数据访问接口
type KeyRepository interface {
	Save(userID int, publicKey string) error
	Get(userID int) (string, error)
	Exists(userID int) (bool, error)
}

type keyRepository struct {
	db *sql.DB
}

// NewKeyRepository 创建密钥仓库实例
func NewKeyRepository(db *sql.DB) KeyRepository {
	return &keyRepository{db: db}
}

func (r *keyRepository) Save(userID int, publicKey string) error {
	_, err := r.db.Exec(
		`INSERT INTO public_keys (user_id, public_key) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET public_key = $2`,
		userID, publicKey,
	)
	return err
}

func (r *keyRepository) Get(userID int) (string, error) {
	var publicKey string
	err := r.db.QueryRow(
		"SELECT public_key FROM public_keys WHERE user_id = $1",
		userID,
	).Scan(&publicKey)

	if err == sql.ErrNoRows {
		return "", errors.New("public key not found")
	}
	if err != nil {
		return "", err
	}

	return publicKey, nil
}

func (r *keyRepository) Exists(userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM public_keys WHERE user_id = $1)",
		userID,
	).Scan(&exists)

	return exists, err
}
