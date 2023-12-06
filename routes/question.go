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
	// TODO: Add routeController.userMiddleware.ExtractUser(),
	router := rg.Group("/question/:group", routeController.questionMiddleware.ExtractGroup())

	router.GET("/all", routeController.questionController.AllQuestions)
	router.GET("/random", routeController.questionController.RandomQuestions)

	// TODO: User must be the owner of the group in order to call these endpoints
	router.POST("/", routeController.questionController.AddQuestion)
	router.PUT("/", routeController.questionController.EditQuestion)
	router.DELETE("/", routeController.questionController.DeleteQuestion)
	router.GET("/", routeController.questionController.GetQuestionById)
}
