package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPhone       = errors.New("invalid phone")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrPhoneRequired      = errors.New("phone is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrPasetoSecretKeyNotFound        = errors.New("paseto secret key not found")
	ErrPasetoRefreshSecretKeyNotFound = errors.New("paseto refresh secret key not found")
)
