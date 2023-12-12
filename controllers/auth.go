package controllers

import (
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	userService  *services.UserService
	tokenService *services.TokenService
	config       *cfg.Config
}

func NewAuthController(userService *services.UserService, tokenService *services.TokenService, config *cfg.Config) *AuthController {
	return &AuthController{
		userService:  userService,
		tokenService: tokenService,
		config:       config,
	}
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	var user models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	newUser, err := controller.userService.CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			utils.ResponseWithStatusMessage(ctx, http.StatusConflict, models.StatusFail, err.Error(), nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusBadGateway, models.StatusFail, err.Error(), nil)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(newUser.ID.Hex())
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(newUser.ID.Hex())
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithStatusMessage(ctx, http.StatusCreated, models.StatusSuccess, "", map[string]any{
		"data": gin.H{"user": models.FilteredResponse(newUser)},
	})
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
	var credentials *models.SignInInput
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	user, err := controller.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, "Invalid email or password", nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, "Invalid email or password", nil)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"access_token": access_token,
	})

}

func (controller *AuthController) SignOut(ctx *gin.Context) {
	refreshToken, _ := ctx.Cookie("refresh_token")

	if err := controller.tokenService.DeleteRefreshToken(refreshToken); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (controller *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusForbidden, models.StatusFail, message, nil)
		return
	}

	// Check if refresh token for the current user exists
	sub, err := controller.tokenService.GetSub(refreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "refresh token doesn't exist") {
			utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, err.Error(), nil)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(sub)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	if err := controller.tokenService.DeleteRefreshToken(refreshToken); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(sub)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"access_token": access_token,
	})
}
