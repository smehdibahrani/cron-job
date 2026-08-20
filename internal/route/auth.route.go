package route

import (
	"cron_job/internal/handler/auth"
	"cron_job/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(api *gin.RouterGroup) {
	route := api.Group("/auth")
	authHandler := auth.NewAuthHandler()

	route.POST("/login", authHandler.Login)
	route.POST("/register", authHandler.Register)

	route.Use(middleware.JWTAuthMiddleware(false))
	route.POST("/refresh-token", authHandler.RefreshToken)

}
