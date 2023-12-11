package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
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
	question := models.Question{
		Group: ctx.GetString("questionGroup"),
	}

	if err := ctx.ShouldBindJSON(&question); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	err := controller.questionService.AddQuestion(question)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, err.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}

func (controller *QuestionController) UpdateQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	var input models.UpdateQuestionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	updatedQuestion := models.Question{
		Group:        input.Group,
		Text:         input.Text,
		Answer:       input.Answer,
		AnswerSource: input.AnswerSource,
	}

	err := controller.questionService.UpdateQuestion(id, updatedQuestion)
	if err != nil {
		fmt.Println(err)
	}

	response(ctx, fmt.Sprintf("edit the #%v word", id))
}

func (controller *QuestionController) DeleteQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	fmt.Println(ctx.GetString("questionGroup"))

	err := controller.questionService.DeleteQuestion(id)

	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, err.Error(), nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, err.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
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
