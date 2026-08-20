package middleware

import (
	"cron_job/internal/usecase"
	"cron_job/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(validateAccessToken bool) gin.HandlerFunc {
	keyAuth := "Authorization"
	if !validateAccessToken {
		keyAuth = "refresh-token"
	}
	return func(c *gin.Context) {
		userUseCase := usecase.NewUserUseCase()
		authHeader := c.GetHeader(keyAuth)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"exception": keyAuth + " header required"})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" || len(bearerToken[1]) < 30 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"exception": "Invalid token format"})
			return
		}

		tokenString := bearerToken[1]
		userID, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"exception": "Invalid token"})
			return
		}
		user, errGet := userUseCase.GetById(userID)
		if errGet != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"exception": "Invalid user"})
			return
		}

		if !user.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"exception": "user is not active"})
			return
		}
		c.Set("userID", userID)
		c.Set("user", user)
		c.Next()
	}
}
