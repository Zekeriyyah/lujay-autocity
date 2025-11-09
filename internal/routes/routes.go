package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/internal/config"
	"github.com/zekeriyyah/lujay-autocity/internal/database"
	"github.com/zekeriyyah/lujay-autocity/internal/handlers"
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
	// r.POST("/auth/login", authHandler.Login)

	// // Protected routes (require authentication)
	// protected := r.Group("/") // Create a group for protected routes
	// protected.Use(middleware.Auth(cfg.JWTSecret)) // Apply Auth middleware to the group

	// {
	// 	// User profile route (accessible by any authenticated user)
	// 	protected.GET("/auth/profile", authHandler.GetProfile)

	// 	// Seller-specific routes (require Seller role)
	// 	sellerRoutes := protected.Group("/") // Create a subgroup for seller routes
	// 	sellerRoutes.Use(middleware.RBAC(types.RoleSeller)) // Apply Seller RBAC middleware
	// 	{
	// 		sellerRoutes.POST("/listings", listingHandler.CreateListing)
	// 		sellerRoutes.GET("/listings/my", listingHandler.GetListingsBySeller) // Adjust handler if needed
	// 		// sellerRoutes.GET("/listings/my/:id", listingHandler.GetListingByID) // Example: Get own listing by ID
	// 		sellerRoutes.PUT("/listings/my/:id", listingHandler.UpdateListing)  // Adjust handler if needed
	// 		sellerRoutes.DELETE("/listings/my/:id", listingHandler.DeleteListing) // Adjust handler if needed
	// 	}

	// 	// Buyer-specific routes (require Buyer role)
	// 	buyerRoutes := protected.Group("/") // Create a subgroup for buyer routes
	// 	buyerRoutes.Use(middleware.RBAC(types.RoleBuyer)) // Apply Buyer RBAC middleware
	// 	{
	// 		buyerRoutes.GET("/listings", listingHandler.GetAllListings) // Only active listings
	// 		buyerRoutes.GET("/listings/:id", listingHandler.GetListingByID) // Only active listings
	// 		// r.POST("/transactions", transactionHandler.InitiatePurchase) // Example
	// 		// r.GET("/transactions/my-purchases", transactionHandler.GetMyPurchases) // Example
	// 	}

	// 	// Admin-specific routes (require Admin role)
	// 	adminRoutes := protected.Group("/") // Create a subgroup for admin routes
	// 	adminRoutes.Use(middleware.RBAC(types.RoleAdmin)) // Apply Admin RBAC middleware
	// 	{
	// 		// r.GET("/users", userHandler.GetAllUsers) // Example
	// 		// r.GET("/users/:id", userHandler.GetUserByID) // Example
	// 		// r.PUT("/users/:id", userHandler.UpdateUser) // Example
	// 		// r.DELETE("/users/:id", userHandler.DeleteUser) // Example
	// 		adminRoutes.GET("/listings", listingHandler.GetAllListings) // All listings, including non-active
	// 		adminRoutes.PUT("/listings/:id", listingHandler.UpdateListing) // For approving/rejecting
	// 		adminRoutes.DELETE("/listings/:id", listingHandler.DeleteListing) // Admin override
	// 		// r.GET("/inspections", inspectionHandler.GetAllInspections) // Example
	// 		// r.PUT("/inspections/:id", inspectionHandler.UpdateInspection) // Example
	// 	}

	// 	// Routes accessible by multiple roles (e.g., Seller and Admin) - Define as needed
	// 	// multiRoleRoutes := protected.Group("/")
	// 	// multiRoleRoutes.Use(middleware.RBAC(types.RoleSeller, types.RoleAdmin))
	// 	// {
	// 	//     // r.GET("/inspections/my-listings/:listingId", inspectionHandler.GetInspectionForMyListing) // Example
	// 	// }
	// }

	return r
}