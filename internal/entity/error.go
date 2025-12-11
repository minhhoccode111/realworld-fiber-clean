package entity

import "errors"

var (
	ErrNoRows                  = errors.New("record not found")
	ErrForbidden               = errors.New("action forbidden")
	ErrNoEffect                = errors.New("zero rows effected")
	ErrConflict                = errors.New("data conflict")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
