package db

import (
	"database/sql"
	"errors"
)

func SavePublicKey(db *sql.DB, userID int, publicKey string) error {
	_, err := db.Exec(
		`INSERT INTO public_keys (user_id, public_key) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET public_key = $2`,
		userID, publicKey,
	)
	return err
}

func GetPublicKey(db *sql.DB, userID int) (string, error) {
	var publicKey string
	err := db.QueryRow(
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

func PublicKeyExists(db *sql.DB, userID int) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM public_keys WHERE user_id = $1)",
		userID,
	).Scan(&exists)

	return exists, err
}
