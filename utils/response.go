package utils

import (
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/gin-gonic/gin"
)

func ResponseWithStatusMessage(ctx *gin.Context, code int, status models.Status, message string) {
	message = strings.TrimSpace(message)
	statusMessage := models.StatusMessage(status)

	h := gin.H{
		"status": statusMessage,
	}
	if message != "" {
		h["message"] = message
	}

	switch status {
	case models.StatusFail:
		ctx.AbortWithStatusJSON(code, h)
	case models.StatusSuccess:
		ctx.JSON(code, h)
	}
}
