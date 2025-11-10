package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"gorm.io/gorm"
)


type InspectionService struct {
	repo *repositories.InspectionRepository
	listingRepo *repositories.ListingRepository
}

func NewInspectionService(inspectionRepo *repositories.InspectionRepository, listingRepo *repositories.ListingRepository) *InspectionService {
	return &InspectionService{
		repo: inspectionRepo,
		listingRepo: listingRepo,
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


// GetInspectionByID retrieves a specific inspection by id
func (s *InspectionService) GetInspectionByID(id string) (*models.Inspection, error) {

	inspection, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection with ID '%s': %w", id, err)
	}
	return inspection, nil
}

// UpdateInspectionStatus updates the status of an inspection and triggers the associated listing status update.
func (s *InspectionService) UpdateInspectionStatus(id string, newStatus models.InspectionStatus, adminID uuid.UUID) error {
	
	// get existing inspection by id
	existingInspection, err := s.repo.GetByID(id)
	if err != nil {
		return err 
	}

	// retrieve the associated listing to update its status
	associatedListing, err := s.listingRepo.GetByID(existingInspection.ListingID.String())
	if err != nil {
		return fmt.Errorf("failed to get associated listing for inspection %s: %w", id, err)
	}

	
	inspectionToUpdate := &models.Inspection{
		ID:     existingInspection.ID, 
		Status: newStatus,             
	}

	// handle listing and inspection update in transaction
	err = s.repo.DB.Transaction(func(tx *gorm.DB) error {
	
		if err := s.repo.UpdateWithTx(tx, inspectionToUpdate); err != nil {
			return fmt.Errorf("failed to update inspection status within transaction: %w", err)
		}

		// update listing status according to current inspection status
		var newListingStatus models.ListingStatus
		switch newStatus {
		case models.InspectionStatusApproved:
			newListingStatus = models.ListingStatusActive
		case models.InspectionStatusRejected:
			newListingStatus = models.ListingStatusRejected
		default:
			return fmt.Errorf("cannot map inspection status '%s' to a listing status", newStatus)
		}

		// Prepare the listing object for update
		listingToUpdate := &models.Listing{
			ID:     associatedListing.ID, 
			Status: newListingStatus,    
		}

		if err := s.listingRepo.UpdateWithTx(tx, listingToUpdate); err != nil {
			return fmt.Errorf("failed to update associated listing status within transaction: %w", err)
		}

		return nil
	})

	// handle error from transaction if any
	if err != nil {

		return fmt.Errorf("failed to update inspection and associated listing: %w", err)
	}
	return nil
}



// CreateInspection creates a new inspection report and conditionally updates the associated listing's status.

func (s *InspectionService) CreateInspection(inspectionToCreate *models.Inspection, adminID uuid.UUID) (*models.Inspection, error) {
	// check for valid listing id
	if inspectionToCreate.ListingID == uuid.Nil {
		return nil, fmt.Errorf("listing_id is required to create an inspection")
	}

	associatedListing, err := s.listingRepo.GetByID(inspectionToCreate.ListingID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get associated listing for inspection: %w", err)
	}

	// switch to know if listing status should be updated or not
	var newListingStatus models.ListingStatus
	statusUpdateRequired := false

	if associatedListing.Status == models.ListingStatusRejected {
		newListingStatus = models.ListingStatusPending
		statusUpdateRequired = true
	} else {
		statusUpdateRequired = false
	}


	// Use a database transaction to handle inspection creation and listing update together
	err = s.repo.DB.Transaction(func(tx *gorm.DB) error {
	
		if err := s.repo.CreateWithTx(tx, inspectionToCreate); err != nil {
			return fmt.Errorf("failed to create inspection within transaction: %w", err)
		}

		if statusUpdateRequired {
		
			listingToUpdate := &models.Listing{
				ID:     associatedListing.ID, 
				Status: newListingStatus,     
			}

			// Update the listing status 
			if err := s.listingRepo.UpdateWithTx(tx, listingToUpdate); err != nil {
				return fmt.Errorf("failed to update associated listing status within transaction: %w", err)
			}
		}

		return nil
	})

	// 
	if err != nil {
		return nil, fmt.Errorf("failed to create inspection and potentially update associated listing: %w", err)
	}

	return inspectionToCreate, nil
}