package controllers

import (
	"fmt"
	"net/http"

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

func (controller *QuestionController) AddQuestion(ctx *gin.Context) {
	response(ctx, "add word")
}

func (controller *QuestionController) UpdateQuestion(ctx *gin.Context) {
	id := ctx.Query("id")
	response(ctx, fmt.Sprintf("edit the #%v word", id))
}

func (controller *QuestionController) DeleteQuestion(ctx *gin.Context) {
	id := ctx.Query("id")
	response(ctx, fmt.Sprintf("delete the #%v word", id))
}

func (controller *QuestionController) GetQuestionById(ctx *gin.Context) {
	id := ctx.Query("id")
	response(ctx, fmt.Sprintf("get the #%v word", id))
}

func (controller *QuestionController) RandomQuestions(ctx *gin.Context) {
	count := ctx.Query("count")
	response(ctx, fmt.Sprintf("get %v random questions", count))
}

func (controller *QuestionController) AllQuestions(ctx *gin.Context) {
	response(ctx, "get all questions")
}

func response(ctx *gin.Context, message string) {
	questionGroup := ctx.GetString("questionGroup")

	ctx.JSON(http.StatusFound, gin.H{
		"group":   questionGroup,
		"message": message,
	})
}
