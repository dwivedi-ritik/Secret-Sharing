package auth

import "errors"

var (
	ErrLoginCredsInvalid error = errors.New("login credentials are invalid")
)
