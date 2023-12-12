package routes

import (
	"github.com/AkifhanIlgaz/random-question-selector/controllers"
	"github.com/AkifhanIlgaz/random-question-selector/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController     *controllers.UserController
	userMiddleware     *middleware.UserMiddleware
	questionMiddleware *middleware.QuestionMiddleware
}

func NewUserRouteController(userController *controllers.UserController, userMiddleware *middleware.UserMiddleware, questionMiddleware *middleware.QuestionMiddleware) UserRouteController {
	return UserRouteController{
		userController:     userController,
		userMiddleware:     userMiddleware,
		questionMiddleware: questionMiddleware,
	}
}

func (routerController *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.Use(routerController.userMiddleware.ExtractUser())
	
	// ? /assign?id=<id>&group=<group>
	router.POST("/assign", routerController.questionMiddleware.ExtractGroup(), routerController.userMiddleware.IsAdminOfGroup(), routerController.userController.AssignGroup)
}
