package middleware

import (
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
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
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotLoggedIn)
			return
		}

		// ! Access token is expired situations must be handled by the client
		claims, err := middleware.tokenService.ParseAccessToken(accessToken)
		if err != nil {
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			return
		}

		user, err := middleware.userService.FindUserById(claims.Subject)
		if err != nil {
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNoUser)
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
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotLoggedIn)
			return
		}

		user, ok := raw.(*models.User)
		if !ok {
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotLoggedIn)
			return
		}

		if !contains(user.Groups, group) {
			httputil.NewError(ctx, http.StatusUnauthorized, utils.ErrNotAdmin)
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
