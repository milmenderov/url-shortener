package storage

import "errors"

var (
	ErrURLNotFound  = errors.New("url not found")
	ErrUserNotExist = errors.New("user not exist")
)
