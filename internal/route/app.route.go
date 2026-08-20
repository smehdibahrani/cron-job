package route

import (
	"context"
	"cron_job/internal/config"
	"cron_job/pkg/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupAppRoutes(api *gin.RouterGroup) {
	route := api.Group("/app")
	route.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": "0.0.1"})
	})

	route.GET("/health", func(c *gin.Context) {
		// Check PostgresSQL connection
		if err := config.DB.Exec("SELECT 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "exception", "service": "PostgreSQL", "exception": err.Error()})
			return
		}

		// Check Redis connection
		if _, err := redis.Rdb.Ping(context.Background()).Result(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "exception", "service": "Redis", "exception": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"connection to redis": "ok", "connection to postgres": "ok"})
	})
}
