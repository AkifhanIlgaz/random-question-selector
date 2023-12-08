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
	router.Use(routerController.userMiddleware.ExtractUser())

	router.GET("/me", routerController.userController.GetMe)

	// ! Possible Errors
	// The assigner isn't the admin of the group
	// The assignee is already the admin of the group

	// ! Middlewares
	// IsAdminOfTheGroup

	// ? /assign?id=<id>&group=<group>
	router.POST("/assign", routerController.userController.AssignGroup)
}
