package main

import (
	"cron_job/internal/config"
	"cron_job/internal/middleware"
	"cron_job/internal/route"
	"cron_job/internal/scheduler"
	"cron_job/pkg/redis"
	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()
	// database
	config.ConnectDatabase()
	redis.Init()
	// Scheduler
	scheduler.Init()
	go scheduler.StartScheduler()

	isDebugMode := config.GetBool("DEBUG_MODE", false)
	if isDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	publicApi := r.Group("/api")
	route.SetupAppRoutes(publicApi)
	route.SetupAuthRoutes(publicApi)

	privateApi := r.Group("/api")
	privateApi.Use(middleware.JWTAuthMiddleware(true))
	route.SetupUserRoutes(privateApi)
	route.SetupJobRoutes(privateApi)
	route.SetupGroupRoutes(privateApi)

	ADDR := config.GetString("ADDR", "8080")
	err := r.Run(":" + ADDR)
	if err != nil {
		return
	}
}
