package repository

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("entity not found")
	ErrDuplicateKey = errors.New("duplicate key violation")
)
