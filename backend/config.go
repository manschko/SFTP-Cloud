package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config stores the application configuration
type Config struct {
	SFTPHost     string
	SFTPPort     int
	JWTSecret    string
	ServerPort   string
}

var AppConfig Config

// LoadConfig loads configuration from environment variables or sets defaults
func LoadConfig() error {
	// Get the environment variable
	portStr := os.Getenv("SFTP_PORT")

	// Convert the string to an integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		// Handle the error (e.g., log it and exit)
		log.Fatalf("Failed to convert SFTP_PORT to integer: %v", err)
	}
	// Set defaults
	AppConfig = Config{
		SFTPHost:     os.Getenv("SFTP_URL"),
		SFTPPort:     port,
		JWTSecret:    "lakflakfh", // todo Change this in production!
		ServerPort:   "8000",
	}

	// Override with environment variables if set
	if host := os.Getenv("SFTP_HOST"); host != "" {
		AppConfig.SFTPHost = host
	}

	if portStr := os.Getenv("SFTP_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid SFTP_PORT value: %v", err)
		}
		AppConfig.SFTPPort = port
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		AppConfig.JWTSecret = jwtSecret
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		AppConfig.ServerPort = serverPort
	}

	return nil
}