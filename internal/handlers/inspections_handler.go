package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
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
	// Get query parameter for optional status filtering
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
