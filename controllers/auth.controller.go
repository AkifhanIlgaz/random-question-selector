package controllers

import (
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.IAuthService
	userService services.IUserService
}

func NewAuthController(authService services.IAuthService, userService services.IUserService) AuthController {
	return AuthController{
		authService: authService,
		userService: userService,
	}
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	var user models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newUser, err := controller.authService.SignUp(&user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(newUser)}})
}
