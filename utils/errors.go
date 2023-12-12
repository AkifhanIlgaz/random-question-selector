package utils

import "errors"

var (
	ErrNoQuestion         error = errors.New("question doesn't exist")
	ErrAlreadyAdmin       error = errors.New("user is already admin of the group")
	ErrNotAdmin           error = errors.New("user isn't the admin of group")
	ErrNoUser             error = errors.New("user doesn't exist")
	ErrSomethingWentWrong error = errors.New("something went wrong")
	ErrNoRefreshToken     error = errors.New("refresh token doesn't exist")
	ErrInvalidCredentials error = errors.New("wrong email or password")
	ErrInvalidCount       error = errors.New("invalid count")
)
