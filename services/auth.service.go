package services

import "github.com/AkifhanIlgaz/random-question-selector/models"

type AuthService interface {
	SignUp(*models.SignUpInput) (*models.User, error)
	SignIn(*models.SignInInput) (*models.User, error)
}


