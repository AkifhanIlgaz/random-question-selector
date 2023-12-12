package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/swaggo/swag/example/celler/httputil"

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

// ShowAccount godoc
// @Summary      Add new question
// @Description  add new question
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Router       /question [post]
func (controller *QuestionController) AddQuestion(ctx *gin.Context) {
	var question models.Question
	if err := ctx.ShouldBindJSON(&question); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if !isAdmin(user.Groups, question.Group) {
		httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotAdmin)
		return
	}

	err := controller.questionService.AddQuestion(question)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}

func (controller *QuestionController) UpdateQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	var updates map[string]string
	if err := ctx.BindJSON(&updates); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if newGroup, ok := updates["group"]; ok {
		if !isAdmin(user.Groups, newGroup) {
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotAdmin)
			return
		}
	}

	err := controller.questionService.UpdateQuestion(id, updates)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			httputil.NewError(ctx, http.StatusNotFound, err)
			return
		}
		httputil.NewError(ctx, http.StatusUnauthorized, err)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}

func (controller *QuestionController) DeleteQuestion(ctx *gin.Context) {
	id := ctx.Query("id")

	question, err := controller.questionService.GetQuestion(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			httputil.NewError(ctx, http.StatusNotFound, err)
			return
		}
		httputil.NewError(ctx, http.StatusUnauthorized, err)
		return
	}

	user := ctx.MustGet("currentUser").(*models.User)
	if !isAdmin(user.Groups, question.Group) {
		httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotAdmin)
		return
	}

	err = controller.questionService.DeleteQuestion(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			httputil.NewError(ctx, http.StatusNotFound, err)
			return
		}
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}

func (controller *QuestionController) GetQuestionById(ctx *gin.Context) {
	id := ctx.Query("id")

	question, err := controller.questionService.GetQuestion(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoQuestion) {
			httputil.NewError(ctx, http.StatusNotFound, utils.ErrNoQuestion)
			return
		}
		httputil.NewError(ctx, http.StatusInternalServerError, utils.ErrSomethingWentWrong)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": question,
	})
}

func (controller *QuestionController) RandomQuestions(ctx *gin.Context) {
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, utils.ErrInvalidCount)
		return
	}
	group := ctx.Query("group")

	questions, err := controller.questionService.GetRandomQuestionsByGroup(group, count)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, utils.ErrSomethingWentWrong)
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
		httputil.NewError(ctx, http.StatusInternalServerError, utils.ErrSomethingWentWrong)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": questions,
	})
}

func isAdmin(groups []string, group string) bool {
	return slices.Contains(groups, group)
}
