package errors

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrInvalidLink       = errors.New("invalid link")
	ErrLinkAlreadyExists = errors.New("link already exists")

	ErrCodeAlreadyExists = errors.New("code already exists")
	ErrStorageIsEmpty    = errors.New("storage is empty")

	ErrUserIDAlreadyExists = errors.New("user id already exists")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
