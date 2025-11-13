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
	CreateWithTx(tx *gorm.DB, inspection *models.Inspection) error
	GetByID(id string) (*models.Inspection, error) 
	GetByStatus(status models.InspectionStatus) ([]*models.Inspection, error)
	GetByListingID(listingID string) (*models.Inspection, error) 
	Update(inspection *models.Inspection) error
	UpdateWithTx(tx *gorm.DB, inspection *models.Inspection) error
	Delete(id string) error 
}

// InspectionRepository implements the InspectionRepositoryInterface.
type InspectionRepository struct {
	DB *gorm.DB
}

// NewInspectionRepository creates a new instance of InspectionRepository.
func NewInspectionRepository(db *gorm.DB) *InspectionRepository {
	return &InspectionRepository{DB: db}
}

// Create add a new inspection to db
func (r *InspectionRepository) Create(inspection *models.Inspection) error {
	return r.CreateWithTx(r.DB, inspection)
}

// CreateWithTx add a new inspection to db while in database transaction
func (r *InspectionRepository) CreateWithTx(tx *gorm.DB, inspection *models.Inspection) error {
	if err := tx.Create(inspection).Error; err != nil {
		return fmt.Errorf("failed to create inspection: %w", err)
	}
	return nil
}


// GetByStatus retrieves inspections filtered by their status.
func (r *InspectionRepository) GetByStatus(status models.InspectionStatus) ([]*models.InspectionFetchInput, error) {
	inspectionsInput := []*models.InspectionFetchInput{}
	// Preload related data (Listing, Inspector)
	if err := r.DB.Preload("Listing").Preload("Inspector").Where("status = ?", status).Find(&inspectionsInput).Error; err != nil {
		return nil, fmt.Errorf("failed to get inspections by status: %w", err)
	}
	
	return inspectionsInput, nil
}


func (i *InspectionRepository) GetByListingID(listingID string) (*models.Inspection, error) {
	parsedListingID, err := pkg.StringToUUID(listingID)
	if err != nil {
		return nil, err
	}

	inspection := &models.Inspection{}

	if err := i.DB.Preload("Listing").Preload("Inspector").First(inspection, "listing_id = ?", parsedListingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("inspection for listing ID %s not found", listingID)
		}
		return nil, fmt.Errorf("failed to get inspection by listing ID: %w", err)
	}
	return inspection, nil
}

// Update updates an existing inspection.
func (i *InspectionRepository) Update(inspection *models.Inspection) error {
	return i.UpdateWithTx(i.DB, inspection)
}

// UpdateWithTx updates an existing inspection within an ongoing transaction
func(i *InspectionRepository) UpdateWithTx(tx *gorm.DB, inspection *models.Inspection) error {
	
	result := tx.Save(inspection) 
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

	result := i.DB.Delete(&models.Inspection{}, "id = ?", parsedID) 
	if result.Error != nil {
		return fmt.Errorf("failed to delete inspection: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("inspection with id %s not found", id)
	}
	return nil
}

func (r *InspectionRepository) GetByID(id string) (*models.InspectionFetchInput, error) {
	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	rawInspection := &models.InspectionFetchInput{}

	if err := r.DB.Table("inspections").Preload("Listing").Preload("Inspector").First(rawInspection, "id = ?", parsedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("inspection with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get inspection: %w", err)
	}
	return rawInspection, nil
}