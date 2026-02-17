package service

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidToken       = errors.New("invalid or expired token")

	ErrClassNotFound  = errors.New("class not found")
	ErrNotClassOwner  = errors.New("not the owner of this class")
	ErrCodeGeneration = errors.New("failed to generate unique class code")
)
