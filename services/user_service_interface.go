package services

import "github.com/AkifhanIlgaz/random-question-selector/models"

type IUserService interface {
	FindUserById(string) (*models.User, error)
	FindUserByEmail(string) (*models.User, error)
}
