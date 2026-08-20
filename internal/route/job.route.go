package route

import (
	"cron_job/internal/handler/job"
	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(api *gin.RouterGroup) {
	route := api.Group("/jobs")

	jobHandler := job.NewJobHandler()

	route.POST("", jobHandler.Create)
	route.PUT("/:id", jobHandler.Update)
	route.GET("", jobHandler.GetAll)
	route.GET("/:id", jobHandler.GetById)
	route.GET("/:id/logs", jobHandler.GetLogsByJobId)
	route.DELETE("/:id", jobHandler.Delete)
}
