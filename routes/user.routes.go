package routes

import (
	"github.com/AkifhanIlgaz/random-question-selector/controllers"
	"github.com/AkifhanIlgaz/random-question-selector/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController *controllers.UserController
	userMiddleware *middleware.UserMiddleware
}

func NewUserRouteController(userController *controllers.UserController, userMiddleware *middleware.UserMiddleware) UserRouteController {
	return UserRouteController{
		userController: userController,
		userMiddleware: userMiddleware,
	}
}

func (routerController *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")

	router.Use(routerController.userMiddleware.DeserializeUser())
	router.GET("/me", routerController.userController.GetMe)
}
