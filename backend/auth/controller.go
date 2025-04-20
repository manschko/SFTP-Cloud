package auth
import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"manschko.com/cloud-storage/sftp"
)

// LoginRequest represents the login form data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthController handles authentication operations
type AuthController struct {
	JWTSecret string
	SFTPHost  string
	SFTPPort  int
}

// NewAuthController creates a new auth controller
func NewAuthController(jwtSecret string, sftpHost string, sftpPort int) *AuthController {
	return &AuthController{
		JWTSecret: jwtSecret,
		SFTPHost:  sftpHost,
		SFTPPort:  sftpPort,
	}
}

// Login authenticates a user with SFTP credentials and returns a JWT token
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create SFTP connection with provided credentials
	sftpConn := sftp.NewConnection(
		c.SFTPHost,
		c.SFTPPort,
		loginReq.Username,
		loginReq.Password,
	)

	// Test connection to verify credentials
	if err := sftpConn.TestConnection(); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(loginReq.Username, c.JWTSecret, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}