package users

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(username string, password string) error {
	_, err := r.GetUser(username)
	if err == nil {
		return fmt.Errorf("Username already exist %w: %s", err, username)

	}

	const op = "SaveUser"

	stmt, err := r.db.Prepare("INSERT INTO users (username, password) values ($1, $2) RETURNING id")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(username, password)
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}

	return nil

}
func (r *AuthPostgres) GetUser(username string) (string, error) {
	const op = "GetUser"
	stmt, err := r.db.Prepare("SELECT password FROM users WHERE username = $1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resUser string

	err = stmt.QueryRow(username).Scan(&resUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotExist
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUser, nil
}
