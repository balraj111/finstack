package main

import (
	"fmt"
	"log"

	"finstack/services/user-service/internal/handler"
	"finstack/services/user-service/internal/repository"
	"finstack/services/user-service/internal/service"
	"finstack/services/user-service/pkg/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Dependency injection
	repo := repository.NewInMemoryUserRepo()
	authService := service.NewAuthService(repo)
	handler.SetAuthService(authService)

	// Public routes
	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)

	// Protected routes
	authGroup := r.Group("/user")
	authGroup.Use(auth.JWTMiddleware())
	authGroup.GET("/profile", handler.Profile)

	fmt.Println("User Service running on port 8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
