package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
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

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
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
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "User doesn't exist",
		})
		return
	}

	usersGroups := candidate.Groups
	if slices.Contains(usersGroups, group) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "User is already admin of the group",
		})
		return
	}

	usersGroups = append(usersGroups, group)
	_, err = controller.userService.UpdateUserById(id, "groups", usersGroups)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
