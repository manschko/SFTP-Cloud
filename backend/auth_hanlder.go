package main

import (
	"github.com/gin-gonic/gin"
	"manschko.com/cloud-storage/auth"
)

// handleLogin handles user authentication and returns a JWT token
func handleLogin(c *gin.Context) {
	authController := auth.NewAuthController(
		AppConfig.JWTSecret,
		AppConfig.SFTPHost,
		AppConfig.SFTPPort,
	)
	authController.Login(c)
}

// authMiddleware creates middleware for JWT authentication
func authMiddleware() gin.HandlerFunc {
	return auth.Middleware(AppConfig.JWTSecret)
}
