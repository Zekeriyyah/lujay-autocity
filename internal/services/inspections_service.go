package services

import (
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
)


type InspectionService struct {
	repo *repositories.InspectionRepository
}

func NewInspectionService(inspectionRepo *repositories.InspectionRepository) *InspectionService {
	return &InspectionService{
		repo: inspectionRepo,
	}
}

// GetInspectionsByStatus retrieves inspections filtered by their status
func (s *InspectionService) GetInspectionsByStatus(status models.InspectionStatus) ([]*models.Inspection, error) {
	
	inspections, err := s.repo.GetByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to get inspections by status '%s': %w", status, err)
	}
	return inspections, nil
}

