package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
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
	var question models.Question
	if err := ctx.ShouldBindJSON(&question); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if !isAdmin(user.Groups, question.Group) {
		utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, utils.ErrNotAdmin.Error(), nil)
		return
	}

	err := controller.questionService.AddQuestion(question)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, err.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}

// TODO:
func (controller *QuestionController) UpdateQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	var updatedQuestion models.Question
	if err := ctx.ShouldBindJSON(&updatedQuestion); err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, err.Error(), nil)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if !isAdmin(user.Groups, updatedQuestion.Group) {
		utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, utils.ErrNotAdmin.Error(), nil)
		return
	}

	err := controller.questionService.UpdateQuestion(id, updatedQuestion)
	if err != nil {
		fmt.Println(err)
	}

}

func (controller *QuestionController) DeleteQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	question, err := controller.questionService.GetQuestion(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, utils.ErrNoQuestion.Error(), nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, utils.ErrSomethingWentWrong.Error(), nil)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if !isAdmin(user.Groups, question.Group) {
		utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, utils.ErrNotAdmin.Error(), nil)
		return
	}

	err = controller.questionService.DeleteQuestion(id)
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

	question, err := controller.questionService.GetQuestion(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, utils.ErrNoQuestion.Error(), nil)
			return
		}
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, utils.ErrSomethingWentWrong.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": question,
	})
}

func (controller *QuestionController) RandomQuestions(ctx *gin.Context) {
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusBadRequest, models.StatusFail, "Please give valid count", nil)
		return
	}
	group := ctx.Query("group")

	questions, err := controller.questionService.GetRandomQuestionsByGroup(group, count)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, utils.ErrSomethingWentWrong.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": questions,
	})
}

func (controller *QuestionController) AllQuestions(ctx *gin.Context) {
	group := ctx.Query("group")

	questions, err := controller.questionService.GetAllQuestionsOfGroup(group)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, utils.ErrSomethingWentWrong.Error(), nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": questions,
	})
}

func response(ctx *gin.Context, message string) {
	questionGroup := ctx.GetString("questionGroup")

	ctx.JSON(http.StatusFound, gin.H{
		"group":   questionGroup,
		"message": message,
	})
}

func isAdmin(groups []string, group string) bool {
	return slices.Contains(groups, group)
}
