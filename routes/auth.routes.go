package routes

import (
	"github.com/AkifhanIlgaz/random-question-selector/controllers"
	"github.com/AkifhanIlgaz/random-question-selector/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
	userMiddleware middleware.UserMiddleware
}

func NewAuthRouteController(authController controllers.AuthController, userMiddleware middleware.UserMiddleware) AuthRouteController {
	return AuthRouteController{
		authController: authController,
		userMiddleware: userMiddleware,
	}
}

func (routeController *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", routeController.authController.SignUp)
	router.POST("/login", routeController.authController.SignIn)
	router.GET("/logout", routeController.userMiddleware.DeserializeUser(), routeController.authController.SignOut)
	router.GET("/refresh", routeController.authController.RefreshAccessToken)
}
