package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func BuildReponseStatusMessage(status string, message string) gin.H {
	message = strings.TrimSpace(message)
	if message == "" {
		return gin.H{
			"status": status,
		}
	} else {
		return gin.H{
			"status":  status,
			"message": message,
		}
	}
}
