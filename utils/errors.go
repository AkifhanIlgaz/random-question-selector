package utils

import "errors"

var (
	ErrNoQuestion         error = errors.New("Question doesn't exist")
	ErrNotAdmin           error = errors.New("User isn't the admin of group")
	ErrSomethingWentWrong error = errors.New("Something went wrong")
)
