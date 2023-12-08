package controllers

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/gin-gonic/gin"
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

	fmt.Println("Group", group)
	fmt.Println("Id", id)
}
