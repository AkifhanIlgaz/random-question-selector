package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
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
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	newUser, err := controller.userService.CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			httputil.NewError(ctx, http.StatusConflict, err)
			return
		}
		httputil.NewError(ctx, http.StatusBadGateway, err)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(newUser.ID.Hex())
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(newUser.ID.Hex())
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithSuccess(ctx, gin.H{"user": models.FilteredResponse(newUser)})
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
	var credentials *models.SignInInput
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			httputil.NewError(ctx, http.StatusBadRequest, utils.ErrInvalidCredentials)
			return
		}
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, utils.ErrInvalidCredentials)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithSuccess(ctx, nil)
}

func (controller *AuthController) SignOut(ctx *gin.Context) {
	refreshToken, _ := ctx.Cookie("refresh_token")

	if err := controller.tokenService.DeleteRefreshToken(refreshToken); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	utils.ResponseWithSuccess(ctx, nil)
}

func (controller *AuthController) RefreshAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, utils.ErrNoRefreshToken)
		return
	}

	// Check if refresh token for the current user exists
	sub, err := controller.tokenService.GetSub(refreshToken)
	if err != nil {
		if errors.Is(err, utils.ErrNoRefreshToken) {
			httputil.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	access_token, err := controller.tokenService.GenerateAccessToken(sub)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := controller.tokenService.DeleteRefreshToken(refreshToken); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	refresh_token, err := controller.tokenService.GenerateRefreshToken(sub)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie("access_token", access_token, controller.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, controller.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", controller.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.ResponseWithSuccess(ctx, nil)
}
