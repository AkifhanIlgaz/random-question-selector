package utils

import "errors"

var (
	ErrNoQuestion         error = errors.New("question doesn't exist")
	ErrNotAdmin           error = errors.New("user isn't the admin of group")
	ErrSomethingWentWrong error = errors.New("something went wrong")
	ErrNoRefreshToken     error = errors.New("refresh token doesn't exist")
)
