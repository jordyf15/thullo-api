package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/utils"
)

type AuthMiddleware struct {
	usecase token.Usecase
}

func NewAuthMiddleware(usecase token.Usecase) *AuthMiddleware {
	return &AuthMiddleware{usecase: usecase}
}

func (middleware *AuthMiddleware) AuthenticateJWT(c *gin.Context) {
	noAuth := map[string][]string{
		"POST":   {"/login", "/login/google", "/login/github", "/tokens/refresh", "/register"},
		"GET":    {"/_health"},
		"DELETE": {"/tokens/remove"},
	}

	requestPath := c.FullPath()

	for _, value := range noAuth[c.Request.Method] {
		if value == requestPath {
			c.Next()
			return
		}
	}

	var response map[string]interface{}
	tokenHeader := c.Request.Header.Get("Authorization")

	if tokenHeader == "" {
		response = utils.Message(false, "Missing auth token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 || splitted[0] != "Bearer" {
		response = utils.Message(false, "Invalid/Malformed auth token")
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	tokenPart := splitted[1]
	tk := &models.AccessToken{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_PASSWORD")), nil
	})

	if err != nil {
		if tk.ExpiresAt < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusForbidden, custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrAccessTokenExpired}})
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrMalformedAccessToken}})
		}
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusForbidden,
			custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrInvalidAccessToken}})
		return
	}

	err = middleware.usecase.Use(tk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden,
			custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrInvalidAccessToken}})
		return
	}

	c.Set("current_user_id", tk.UserID)
	c.Next()
}
