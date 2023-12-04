package services

import "github.com/AkifhanIlgaz/random-question-selector/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.User, error)
	SignInUser(*models.SignInInput) (*models.User, error)
}
