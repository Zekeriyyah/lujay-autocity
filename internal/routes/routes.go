package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/internal/config"
	"github.com/zekeriyyah/lujay-autocity/internal/database"
	"github.com/zekeriyyah/lujay-autocity/internal/handlers"
	"github.com/zekeriyyah/lujay-autocity/internal/middleware"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
)

func SetupRouter(r *gin.Engine, cfg *config.Config) *gin.Engine {

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)


	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	// Initialize Handler
	authHandler := handlers.NewAuthHandler(authService)




	//------------------ENDPOINTS------------------------

	// health-check
	r.GET("/ping", func(c *gin.Context){
		c.JSON(200, gin.H{"health_check": "pong"})
	})

	// Public routes (no auth required)
	r.POST("/auth/register", authHandler.Register) 
	r.POST("/auth/login", authHandler.Login)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) 

	{
		// User profile route (accessible by any authenticated user)
		protected.GET("/auth/profile", authHandler.GetProfile)

	return r
	}
}