package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"url-shortener/internal/storage"
)

type Storage struct {
	Db *sql.DB
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {

	const op = "storage.postgres.SaveURL"

	stmt, err := s.Db.Prepare("INSERT INTO url(url, alias) VALUES($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.Db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	_, err := s.GetURL(alias)
	if err != nil {
		return fmt.Errorf("Alias not found %w: %s", err, alias)
	}

	const op = "storage.postgres.DeleteURL"

	stmt, err := s.Db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateURL(newUrl string, alias string) error {

	const op = "storage.postgres.UpdateURL"

	stmt, err := s.Db.Prepare("UPDATE url SET url = $1 WHERE alias = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(newUrl, alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
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
