package routes

import (
	"github.com/AkifhanIlgaz/random-question-selector/controllers"
	"github.com/AkifhanIlgaz/random-question-selector/middleware"
	"github.com/gin-gonic/gin"
)

type QuestionRouteController struct {
	questionController *controllers.QuestionController
	userMiddleware     *middleware.UserMiddleware
	questionMiddleware *middleware.QuestionMiddleware
}

func NewQuestionRouteController(questionController *controllers.QuestionController, userMiddleware *middleware.UserMiddleware, questionMiddleware *middleware.QuestionMiddleware) *QuestionRouteController {
	return &QuestionRouteController{
		questionController: questionController,
		userMiddleware:     userMiddleware,
		questionMiddleware: questionMiddleware,
	}
}

func (routeController *QuestionRouteController) QuestionRoute(rg *gin.RouterGroup) {
	router := rg.Group("/question", routeController.userMiddleware.ExtractUser())

	{
		// ?group=<group>
		router.GET("/all", routeController.questionController.AllQuestions)
		router.GET("/random", routeController.questionController.RandomQuestions)
	}

	adminRoute := router.Group("/")
	{
		adminRoute.POST("", routeController.questionController.AddQuestion)
		adminRoute.PUT("", routeController.questionController.UpdateQuestion)
		adminRoute.DELETE("", routeController.questionController.DeleteQuestion)
		// ? Should be the admin
		adminRoute.GET("", routeController.questionController.GetQuestionById)
	}
}
