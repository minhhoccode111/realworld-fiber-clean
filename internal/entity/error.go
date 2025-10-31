package entity

import "errors"

var ErrNoRows = errors.New("record not found")
var ErrForbidden = errors.New("action forbidden")
