package middleware

import (
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
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
			utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, "You are not logged in", nil)
			return
		}

		// ! Access token is expired situations must be handled by the client
		claims, err := middleware.tokenService.ParseAccessToken(accessToken)
		if err != nil {
			utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, err.Error(), nil)
			return
		}

		user, err := middleware.userService.FindUserById(claims.Subject)
		if err != nil {
			utils.ResponseWithStatusMessage(ctx, http.StatusNotFound, models.StatusFail, "User does not exist", nil)
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

func (middleware *UserMiddleware) IsAdminOfGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		group := ctx.GetString("questionGroup")

		raw, ok := ctx.Get("currentUser")
		if !ok {
			utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, "You are not logged in", nil)
			return
		}

		user, ok := raw.(*models.User)
		if !ok {
			utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, "You are not logged in", nil)
			return
		}

		if !contains(user.Groups, group) {
			utils.ResponseWithStatusMessage(ctx, http.StatusUnauthorized, models.StatusFail, "You are not admin of this group", nil)
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
