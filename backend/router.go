package main

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API endpoints
func SetupRoutes(router *gin.Engine) {
	// Auth routes
	router.POST("/api/login", handleLogin)

	// Protected routes
	authorized := router.Group("/api")
	authorized.Use(authMiddleware())
	{
		// File operations
		authorized.GET("/files/*path", listFiles)
		authorized.GET("/download/*path", downloadFile)
		authorized.POST("/upload/*path", uploadFile)
		authorized.DELETE("/files/*path", deleteFile)
		authorized.PUT("/move", moveFile)
		authorized.PUT("/rename", renameFile)
		authorized.POST("/mkdir/*path", createDirectory)
	}
}
