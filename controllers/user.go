package controllers

import (
	"errors"
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
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

func (controller *UserController) AssignGroup(ctx *gin.Context) {
	group := ctx.Query("group")
	id := ctx.Query("id")

	candidate, err := controller.userService.FindUserById(id)
	if err != nil {
		if errors.Is(err, utils.ErrNoUser) {
			httputil.NewError(ctx, http.StatusNotFound, utils.ErrNoUser)
		}
		return
	}

	usersGroups := candidate.Groups
	if slices.Contains(usersGroups, group) {
		httputil.NewError(ctx, http.StatusNotFound, utils.ErrAlreadyAdmin)
		return
	}

	usersGroups = append(usersGroups, group)
	_, err = controller.userService.UpdateUserById(id, "groups", usersGroups)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, utils.ErrSomethingWentWrong)
		return
	}

	utils.ResponseWithSuccess(ctx, nil)
}
