package controllers

import (
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	questionService *services.QuestionService
}

func NewQuestionController(questionService *services.QuestionService) *QuestionController {
	return &QuestionController{
		questionService: questionService,
	}
}

func (controller *QuestionController) AddQuestion(ctx *gin.Context) {}

func (controller *QuestionController) EditQuestion(ctx *gin.Context) {}

func (controller *QuestionController) DeleteQuestion(ctx *gin.Context) {}

func (controller *QuestionController) GetQuestionById(ctx *gin.Context) {}

func (controller *QuestionController) RandomQuestions(ctx *gin.Context) {}

func (controller *QuestionController) AllQuestions(ctx *gin.Context) {}
