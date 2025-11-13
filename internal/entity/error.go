package entity

import "errors"

var (
	ErrNoRows        = errors.New("record not found")
	ErrForbidden     = errors.New("action forbidden")
	ZeroRowsAffected = errors.New("zero rows effected")
)
