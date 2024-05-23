package main

import (
	"log"
	"vcs_server/cron"
	"vcs_server/handlers"
	"vcs_server/helper"
	"vcs_server/middleware"

	_ "vcs_server/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	helper.Migration()
	go cron.Cron_SendReport()
	go cron.Cron_healcheck()
	go cron.Cron_Invalid_Cache()
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	queryLogger := &lumberjack.Logger{
		Filename:   "./logs/loquery.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}
	queryLog := log.New(queryLogger, "", log.LstdFlags)

	router.POST("/login", middleware.LoggingMiddleware(queryLog), handlers.Login)
	apiRouters := router.Group("/api", middleware.LoggingMiddleware(queryLog), middleware.JWTAuthMiddleware())
	{

		apiRouters.POST("/server", middleware.AdminOnlyMiddleware(), handlers.CreateNewServer)

		apiRouters.PUT("/server/:id", middleware.AdminOnlyMiddleware(), handlers.UpdateServer)

		apiRouters.DELETE("/server/:id", middleware.AdminOnlyMiddleware(), handlers.DeleteServer)

		apiRouters.GET("/servers", handlers.GetOrExportServers)

		apiRouters.POST("/servers", middleware.AdminOnlyMiddleware(), handlers.ImportServers)

		apiRouters.GET("/server/:id", handlers.GetServerById)
		apiRouters.GET("/servers/report", handlers.ReportServer)

		apiRouters.POST("/user", middleware.AdminOnlyMiddleware(), handlers.AddUser)
	}

	router.Run(":8080")
}
