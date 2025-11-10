package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

