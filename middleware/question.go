package middleware

import (
	"github.com/gin-gonic/gin"
)

type QuestionMiddleware struct {
}

func NewQuestionMiddleware() *QuestionMiddleware {
	return &QuestionMiddleware{}
}

func (middleware *QuestionMiddleware) ExtractGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		questionGroup := ctx.Param("group")
		ctx.Set("questionGroup", questionGroup)
	}
}
