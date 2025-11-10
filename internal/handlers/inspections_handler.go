package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/internal/middleware"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
	"gorm.io/gorm"
)


type InspectionHandler struct {
	service   *services.InspectionService
	validator *validator.Validate 
}

type CreateInspectionInput struct {
	ListingID       uuid.UUID           `json:"listing_id" validate:"required,uuid"`
	InspectorID     uuid.UUID           `json:"inspector_id" validate:"required,uuid"`
	InspectionDate  *time.Time          `json:"inspection_date,omitempty"` // Optional: Date of inspection
	ConditionRating *int                `json:"condition_rating,omitempty" validate:"omitempty,min=1,max=10"`
	Findings        map[string]any      `json:"findings,omitempty"` 
	ReportURL       string              `json:"report_url,omitempty" validate:"omitempty,url"` 
	Status          models.InspectionStatus `json:"status" validate:"required,oneof=pending approved rejected"` 
}


func NewInspectionHandler(service *services.InspectionService) *InspectionHandler {
	return &InspectionHandler{
		service:   service,
		validator: validator.New(), 
	}
}

// GetInspections handles GET /inspections (admin only)
func (h *InspectionHandler) GetInspections(c *gin.Context) {
	
	statusParam := c.Query("status")
	var status models.InspectionStatus

	if statusParam != "" {
		switch models.InspectionStatus(statusParam) {
		case models.InspectionStatusPending, models.InspectionStatusApproved, models.InspectionStatusRejected:
			status = models.InspectionStatus(statusParam)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid status parameter: %s", statusParam)})
			return
		}
	} else {
		// Default to 'pending' inspections if no status is specified
		status = models.InspectionStatusPending
	}

	inspections, err := h.service.GetInspectionsByStatus(status)
	if err != nil {
		log.Printf("Error getting inspections by status '%s': %v", status, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inspections"})
		return
	}

	c.JSON(http.StatusOK, inspections)
}


func (h *InspectionHandler) GetInspectionByID(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection ID format"})
		return
	}

	// invoke service layer to perform the get logic
	inspection, err := h.service.GetInspectionByID(idStr)
	if err != nil {
		log.Printf("error getting inspection %s: %v", idStr, err)
		if err.Error() == fmt.Sprintf("inspection with id %s not found", idStr) {
			c.JSON(http.StatusNotFound, gin.H{"error": "inspection not found"})
			return
		}
	
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get inspection"})
		return
	}

	c.JSON(http.StatusOK, inspection)
}


// UpdateInspection handles PUT /inspections/{id} (admin only)
func (h *InspectionHandler) UpdateInspection(c *gin.Context) {

	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inspection ID format"})
		return
	}

	adminID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type UpdateInspectionInput struct {
		Status models.InspectionStatus `json:"status" validate:"required,oneof=pending approved rejected"`
	}

	// Parse JSON request body into input
	var input UpdateInspectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	// Validate the input struct fields using the validator
	if err := h.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation error: %v", err)})
		return
	}

	
	if err := h.service.UpdateInspectionStatus(idStr, input.Status, adminID); err != nil {
		log.Printf("Error updating inspection %s: %v", idStr, err)
		
		if err.Error() == fmt.Errorf("inspection with id %s not found", idStr).Error() ||
		   err.Error() == fmt.Errorf("failed to get associated listing for inspection %s: %w", idStr, gorm.ErrRecordNotFound).Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "inspection or associated listing not found"})
			return
		}

		// Generic error for other failures (e.g., database transaction issues)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inspection"})
		return
	}

	updatedInspection, err := h.service.GetInspectionByID(idStr)
	if err != nil {
		log.Printf("Warning: Inspection %s updated successfully, but failed to fetch updated details: %v", idStr, err)
		c.JSON(http.StatusOK, gin.H{"message": "Inspection status updated successfully"})
		return
	}

	c.JSON(http.StatusOK, updatedInspection)
}



// CreateInspection handles POST /inspections (admin only)
func (h *InspectionHandler) CreateInspection(c *gin.Context) {
	// get id from context
	adminID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse JSON request body into the CreateInspectionInput struct
	var input CreateInspectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation error: %v", err)})
		return
	}

	// map the inspection object for the service call.
	inspectionToCreate := &models.Inspection{
		ListingID:       input.ListingID,
		InspectorID:     input.InspectorID,
		InspectionDate:  time.Time{},
		ConditionRating: 0, 
		Findings:        input.Findings,
		ReportURL:       input.ReportURL,
		Status:          input.Status, 
	}

	// Handle pointer fields from input
	if input.InspectionDate != nil {
		inspectionToCreate.InspectionDate = *input.InspectionDate
	}
	if input.ConditionRating != nil {
		inspectionToCreate.ConditionRating = *input.ConditionRating
	}


	createdInspection, err := h.service.CreateInspection(inspectionToCreate, adminID)
	if err != nil {
		log.Printf("Error creating inspection: %v", err)
		
		if err.Error() == "listing_id is required to create an inspection" ||
		   err.Error() == fmt.Errorf("failed to get associated listing for inspection: %w", gorm.ErrRecordNotFound).Error() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Associated listing not found or invalid listing_id"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create inspection"})
		return
	}

	
	c.JSON(http.StatusCreated, createdInspection)
}