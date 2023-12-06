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
	router := rg.Group("/question/:group", routeController.userMiddleware.ExtractUser(), routeController.questionMiddleware.ExtractGroup())

	router.GET("/all", routeController.questionController.AllQuestions)
	router.GET("/random", routeController.questionController.RandomQuestions)

	router.POST("/add", routeController.questionController.AddQuestion)
	router.PUT("/edit", routeController.questionController.EditQuestion)
	router.DELETE("/", routeController.questionController.DeleteQuestion)
	router.GET("/", routeController.questionController.GetQuestionById)
}
