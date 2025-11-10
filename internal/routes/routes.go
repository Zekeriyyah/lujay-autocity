package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/internal/config"
	"github.com/zekeriyyah/lujay-autocity/internal/database"
	"github.com/zekeriyyah/lujay-autocity/internal/handlers"
	"github.com/zekeriyyah/lujay-autocity/internal/middleware"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
)

func SetupRouter(r *gin.Engine, cfg *config.Config) *gin.Engine {

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	listingRepo := repositories.NewListingRepository(database.DB)
	vehicleRepo := repositories.NewVehicleRepository(database.DB)
	inspectionRepo := repositories.NewInspectionRepository(database.DB)


	// Initialize services
	authService := services.NewAuthService(userRepo)
	listingService := services.NewListingService(listingRepo, vehicleRepo, userRepo)
	inspectionService := services.NewInspectionService(inspectionRepo)

	// Initialize Handler
	authHandler := handlers.NewAuthHandler(authService)
	listingHandler := handlers.NewListingHandler(listingService)
	inspectionHandler := handlers.NewInspectionHandler(inspectionService)



	//------------------ENDPOINTS------------------------

	// health-check
	r.GET("/ping", func(c *gin.Context){
		c.JSON(200, gin.H{"health_check": "pong"})
	})

	// Public routes (no auth required)
	r.POST("/auth/register", authHandler.Register) 
	r.POST("/auth/login", authHandler.Login)
	r.GET("/listings", listingHandler.GetAllActiveListings)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) 

	{
		// User profile route (accessible by any authenticated user)
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Seller-specific routes 
		sellerRoutes := protected.Group("/")
		sellerRoutes.Use(middleware.RBAC(types.RoleSeller))
		{
			sellerRoutes.POST("/listings", listingHandler.CreateListing)
			sellerRoutes.GET("/listings/my", listingHandler.GetListingsBySeller)

		}

		// Buyer-specific routes
		buyerRoutes := protected.Group("/")
		buyerRoutes.Use(middleware.RBAC(types.RoleBuyer))

		// Admin-specific routes
		adminRoutes := protected.Group("/")
		adminRoutes.Use(middleware.RBAC(types.RoleAdmin))

		{
			adminRoutes.GET("/listings/all", listingHandler.GetAllListingsForAdmin)
			adminRoutes.GET("/listings/:id", listingHandler.GetListingByID) 
			adminRoutes.DELETE("/listings/:id", listingHandler.DeleteListing)
			
			adminRoutes.GET("/inspections", inspectionHandler.GetInspections)
			adminRoutes.GET("/inspections/:id", inspectionHandler.GetInspectionByID) 

		}

		// For both seller and admin 
		protectedSellerOrAdmin := protected.Group("/")
		protectedSellerOrAdmin.Use(middleware.RBAC(types.RoleSeller, types.RoleAdmin))
		{
			protectedSellerOrAdmin.PUT("/listings/:id", listingHandler.UpdateListing)
		}

	return r
	}
}