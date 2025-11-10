package handlers

import (
	"fmt"
	"log"
	"net/http"

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
