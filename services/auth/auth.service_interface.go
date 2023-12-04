package services

import "github.com/AkifhanIlgaz/random-question-selector/models"

type IAuthService interface {
	SignUp(*models.SignUpInput) (*models.User, error)
	SignIn(*models.SignInInput) (*models.User, error)
}
