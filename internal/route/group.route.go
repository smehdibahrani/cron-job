package route

import (
	"cron_job/internal/handler/group"
	"github.com/gin-gonic/gin"
)

func SetupGroupRoutes(api *gin.RouterGroup) {
	route := api.Group("/groups")
	groupHandler := group.NewGroupHandler()

	route.POST("/", groupHandler.Create)
	route.PUT("/:id", groupHandler.Update)
	route.GET("/", groupHandler.GetAll)
	route.GET("/:id", groupHandler.GetById)
	route.GET("/:id/jobs", groupHandler.GetJobsByGrpId)
	route.DELETE("/:id", groupHandler.Delete)
}
