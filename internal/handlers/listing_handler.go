package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zekeriyyah/lujay-autocity/internal/middleware"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
)

type ListingHandler struct {
	service *services.ListingService
	validator *validator.Validate
}

func NewListingHandler(service *services.ListingService) *ListingHandler{
	return &ListingHandler{
		service: service,
		validator: validator.New(),
	}
}

// Input structure for creating a listing and its associated vehicle
type CreateListingInput struct {
	Title       string        `json:"title" validate:"required"`
	Description string        `json:"description,omitempty"`
	Price       float64       `json:"price" validate:"required,gt=0"`
	Location    string        `json:"location" validate:"required"`
	
	// Vehicle details
	Vehicle models.Vehicle `json:"vehicle" validate:"required"` // Validate the nested Vehicle struct
}

// CreateListing handles POST /listings (requires Seller role via middleware)

func (h *ListingHandler) CreateListing(c *gin.Context) {
	// Extract authenticated Seller ID from Gin context
	authUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse JSON request body into the combined CreateListingInput struct
	input := &CreateListingInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON", "details": err.Error()})
		return
	}

	// Validate the combined input struct fields using the validator
	if err := h.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation error: %v", err)})
		return
	}

	// Prepare the Listing model from the input
	listing := &models.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Location:    input.Location,
		
		SellerID: authUserID,
	}

	vehicle := &input.Vehicle

	// Create vehicle if not in existence then create listing
	createdListing, err := h.service.CreateListingWithVehicle(listing, vehicle, authUserID)
	if err != nil {
		log.Printf("Error creating listing with vehicle: %v", err)
		
		if err.Error() == "unauthorized: seller ID does not match authenticated user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create listing"})
		return
	}

	c.JSON(http.StatusCreated, createdListing)
}