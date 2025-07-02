package handler

import (
	"net/http"

	"finstack/services/user-service/internal/model"
	"finstack/services/user-service/internal/service"
	"finstack/services/user-service/pkg/auth"

	"github.com/gin-gonic/gin"
)

var authService *service.AuthService

func SetAuthService(s *service.AuthService) {
	authService = s
}

func SignUp(c *gin.Context) {
	var req model.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := authService.SignUp(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user.Email})
}

func Profile(c *gin.Context) {
	email := c.GetString("email")
	userID := c.GetString("userID")
	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
		"email":  email,
	})
}

func Logout(c *gin.Context) {
	token := c.GetString("token")
	exp := c.GetFloat64("exp")
	if token == "" || exp == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token context"})
		return
	}

	auth.AddToken(token, int64(exp))
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
