package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
	"gorm.io/gorm"
)


type ListingService struct {
	repo *repositories.ListingRepository
	vehicleRepo *repositories.VehicleRepository
	userRepo *repositories.UserRepository
}

func NewListingService(repo *repositories.ListingRepository, vehicleRepo *repositories.VehicleRepository, userRepo *repositories.UserRepository) *ListingService {
	return &ListingService{
		repo: repo,
		vehicleRepo: vehicleRepo,
		userRepo: userRepo,
	}
}

// CreateListingWithVehicle perform creation of a new listing for a seller.
func (s *ListingService) CreateListingWithVehicle(listing *models.Listing, vehicle *models.Vehicle, authenticatedUserID uuid.UUID) (*models.Listing, error) {
	// Validate that the authenticated user is the seller attempting to create the listing
	if listing.SellerID != authenticatedUserID {
		return nil, fmt.Errorf("unauthorized: seller ID does not match authenticated user")
	}

	// using GORM's Transaction to handle operations on vehicle and listing in a roll
	var createdListing *models.Listing
	err := s.repo.DB.Transaction(func(tx *gorm.DB) error { 
		
		// Check if vehicle of same vin exists to use it or create new one
		existingVehicle, err := s.vehicleRepo.GetByVINWithTx(tx, vehicle.VIN)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("error checking for existing vehicle: %w", err)
			}
			
			if err := s.vehicleRepo.CreateWithTx(tx, vehicle); err != nil {
				return fmt.Errorf("failed to create vehicle: %w", err)
			}
		
		} else {
			vehicle = existingVehicle 
		}

		// link the listing to the newly created or existing vehicle
		listing.VehicleID = vehicle.ID
		listing.Status = models.ListingStatusPending 

		// create the listing within the same transaction
		if err := s.repo.CreateWithTx(tx, listing); err != nil {
			return fmt.Errorf("failed to create listing: %w", err)
		}

		createdListing = listing
		return nil 
	})

	// handle error from transaction
	if err != nil {
		return nil, err
	}

	return createdListing, nil
}

// GetAllListings retrieves all listings with the status 'active'.
func (s *ListingService) GetAllListings() ([]*models.Listing, error) {
	return s.repo.GetAll()
}

// GetActiveListings retrieves all listings with the status 'active'.
func (s *ListingService) GetActiveListings() ([]*models.Listing, error) {
	return s.repo.GetByStatus(models.ListingStatusActive)
}


// GetListingsByStatus retrieves listings filtered by a specific status
func (s *ListingService) GetListingsByStatus(status models.ListingStatus) ([]*models.Listing, error) {
	return s.repo.GetByStatus(status)
}


// GetListingsBySellerID retrieves all listings for a specific seller (used by Sellers).
func (s *ListingService) GetListingsBySellerID(sellerID string) ([]*models.Listing, error) {
	return s.repo.GetBySellerID(sellerID)
}


// GetListingByID retrieves a specific listing by its ID.
func (s *ListingService) GetListingByID(id string) (*models.Listing, error) {
	return s.repo.GetByID(id)
}


func (s *ListingService) UpdateListing(listingToUpdate *models.Listing, authenticatedUserID uuid.UUID, role types.Role) error {
	// Retrieve the existing listing from the database
	existingListing, err := s.repo.GetByID(listingToUpdate.ID.String())
	if err != nil {
		return err
	}

	// Check authorization based on role
	switch role {
	case types.RoleSeller:
		// Seller can only update their own listing
		if existingListing.SellerID != authenticatedUserID {
			return fmt.Errorf("unauthorized: you can only update your own listings")
		}
		// Seller cannot update sellerID, VehicleID and status
		listingToUpdate.SellerID = existingListing.SellerID 
		listingToUpdate.VehicleID = existingListing.VehicleID
		listingToUpdate.Status = existingListing.Status 

	case types.RoleAdmin:
		// Admin can update any listing including status but not sellerID and VehicleID
		listingToUpdate.SellerID = existingListing.SellerID 
		listingToUpdate.VehicleID = existingListing.VehicleID 
		

	default:
		return fmt.Errorf("unauthorized: invalid role for update operation")
	}

	return s.repo.Update(listingToUpdate)
}

func (s *ListingService) DeleteListing(id string) error {
	return s.repo.Delete(id)
}