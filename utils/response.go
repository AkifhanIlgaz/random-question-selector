package utils

import (
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/gin-gonic/gin"
)

func ResponseWithStatusMessage(ctx *gin.Context, code int, status models.Status, message string, additionalInfo map[string]any) {
	message = strings.TrimSpace(message)
	statusMessage := models.StatusMessage(status)

	h := gin.H{
		"status": statusMessage,
	}
	if message != "" {
		h["message"] = message
	}
	for k, v := range additionalInfo {
		h[k] = v
	}

	switch status {
	case models.StatusFail:
		ctx.AbortWithStatusJSON(code, h)
	case models.StatusSuccess:
		ctx.JSON(code, h)
	}
}
