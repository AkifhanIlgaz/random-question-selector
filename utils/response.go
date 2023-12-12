package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseWithSuccess(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
