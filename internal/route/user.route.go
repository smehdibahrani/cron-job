package route

import (
	"cron_job/internal/handler/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup) {
	route := api.Group("/users")
	userHandler := user.NewUserHandler()

	route.PUT("/", userHandler.Update)
	route.GET("/", userHandler.GetInfo)
	route.POST("/change-password", userHandler.ChangePassword)

}
