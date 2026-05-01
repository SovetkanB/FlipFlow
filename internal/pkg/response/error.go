package response

import "errors"

var (
	ErrNotFound             = errors.New("record not found")
	ErrEmailTaken           = errors.New("email already in use")
	ErrInvalidPassword      = errors.New("email or password is not correct")
	ErrNotValidRefreshToken = errors.New("not valid refresh token")
)
