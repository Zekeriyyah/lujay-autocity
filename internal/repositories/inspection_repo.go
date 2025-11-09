package repositories

import (
	"errors"
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/pkg"
	"gorm.io/gorm"
)

type InspectionRepositoryInterface interface {
	Create(inspection *models.Inspection) error
	GetByID(id string) (*models.Inspection, error) 
	GetByListingID(listingID string) (*models.Inspection, error) 
	Update(inspection *models.Inspection) error
	Delete(id string) error 
}

// InspectionRepository implements the InspectionRepositoryInterface.
type InspectionRepository struct {
	db *gorm.DB
}

// NewInspectionRepository creates a new instance of InspectionRepository.
func NewInspectionRepository(db *gorm.DB) *InspectionRepository {
	return &InspectionRepository{db: db}
}

func (i *InspectionRepository) GetByListingID(listingID string) (*models.Inspection, error) {
	parsedListingID, err := pkg.StringToUUID(listingID)
	if err != nil {
		return nil, err
	}

	inspection := &models.Inspection{}

	if err := i.db.Preload("Listing").Preload("Inspector").First(inspection, "listing_id = ?", parsedListingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("inspection for listing ID %s not found", listingID)
		}
		return nil, fmt.Errorf("failed to get inspection by listing ID: %w", err)
	}
	return inspection, nil
}

// Update updates an existing inspection.
func (i *InspectionRepository) Update(inspection *models.Inspection) error {
	
	result := i.db.Save(inspection) 
	if result.Error != nil {
		return fmt.Errorf("failed to update inspection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("inspection with id %s not found", inspection.ID.String()) 
	}
	return nil
}

// Delete removes an inspection by its UUID ID (passed as a string).
func (i *InspectionRepository) Delete(id string) error {
	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return err
	}

	result := i.db.Delete(&models.Inspection{}, "id = ?", parsedID) 
	if result.Error != nil {
		return fmt.Errorf("failed to delete inspection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("inspection with id %s not found", id)
	}
	return nil
}
