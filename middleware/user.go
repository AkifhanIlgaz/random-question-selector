package middleware

import (
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	userService  *services.UserService
	tokenService *services.TokenService
}

func NewUserMiddleware(userService *services.UserService, tokenService *services.TokenService) *UserMiddleware {
	return &UserMiddleware{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (middleware *UserMiddleware) ExtractUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		// Read access token from cookie or authorizaiton header
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		claims, err := middleware.tokenService.ParseAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.Set("currentUser", services.RedisClaims{
			Subject: claims.Subject,
			Groups:  claims.Groups,
		})
		ctx.Next()
	}
}

func (middleware *UserMiddleware) IsAdminOfGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		group := ctx.GetString("questionGroup")
		raw, ok := ctx.Get("currentUser")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		user, ok := raw.(services.RedisClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		if !contains(user.Groups, group) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not admin of this group"})
			return
		}

		ctx.Next()
	}
}

func contains(groups []string, group string) bool {
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}
