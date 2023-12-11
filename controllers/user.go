package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (controller *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", map[string]any{
		"data": gin.H{
			"user": models.FilteredResponse(currentUser),
		},
	})
}

func (controller *UserController) AssignGroup(ctx *gin.Context) {
	group := ctx.Query("group")
	id := ctx.Query("id")

	candidate, err := controller.userService.FindUserById(id)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, "User doesn't exist", nil)
		return
	}

	usersGroups := candidate.Groups
	if slices.Contains(usersGroups, group) {
		utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, "User is already admin of the group", nil)
		return
	}

	usersGroups = append(usersGroups, group)
	_, err = controller.userService.UpdateUserById(id, "groups", usersGroups)
	if err != nil {
		utils.ResponseWithStatusMessage(ctx, http.StatusInternalServerError, models.StatusFail, "Something went wrong", nil)
		return
	}

	utils.ResponseWithStatusMessage(ctx, http.StatusOK, models.StatusSuccess, "", nil)
}
