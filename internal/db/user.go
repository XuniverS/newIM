package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
}

func CreateUser(db *sql.DB, username, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	err = db.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		username, string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(db *sql.DB, userID int) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, username FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
