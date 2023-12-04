package controllers

import "github.com/AkifhanIlgaz/random-question-selector/services"

type AuthController struct {
	authService services.IAuthService
	userService services.IUserService
}


