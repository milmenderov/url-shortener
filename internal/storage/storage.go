package storage

import "errors"

var (
	ErrURLExists   = errors.New("url exists")
	ErrURLNotFound = errors.New("url not found")
)
