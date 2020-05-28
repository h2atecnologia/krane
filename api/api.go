package api

import (
	"strings"

	"github.com/biensupernice/krane/api/handler"
	"github.com/biensupernice/krane/api/middleware"
	"github.com/biensupernice/krane/internal/deployment/event"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config : server config
type Config struct {
	RestPort string // port to expose rest api
	LogLevel string // release | debug
}

// Start : api server
func Start(cnf Config) {
	gin.SetMode(strings.ToLower(cnf.LogLevel))

	client := gin.New()

	// Middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.POST("/health", handler.HealthHandler)
	client.POST("/auth", handler.Auth)
	client.GET("/login", handler.Login)

	client.GET("/sessions", middleware.AuthSessionMiddleware(), handler.GetSessions)

	client.POST("/deployments", middleware.AuthSessionMiddleware(), handler.CreateSpec)
	client.GET("/deployments", middleware.AuthSessionMiddleware(), handler.GetDeployments)
	client.GET("/deployments/:name", middleware.AuthSessionMiddleware(), handler.GetDeployment)
	client.DELETE("/deployments/:name", middleware.AuthSessionMiddleware(), handler.DeleteDeployment)
	client.POST("/deployments/:name/run", middleware.AuthSessionMiddleware(), handler.RunDeployment)

	// -- SSE -- //
	client.GET("/containers/:containerID/events", middleware.AuthSessionMiddleware(), handler.ContainerEvents)

	// --  Websockets -- //
	client.GET("/deployments/:name/events", handler.WSDeploymentHandler)
	go event.Echo()

	client.Run(":" + cnf.RestPort)
}