package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"gorm.io/gorm"
)


type ListingService struct {
	repo repositories.ListingRepositoryInterface
	vehicleRepo repositories.VehicleRepositoryInterface
	userRepo repositories.UserRepositoryInterface
}

func NewListingService(repo repositories.ListingRepositoryInterface, vehicleRepo repositories.VehicleRepositoryInterface, userRepo repositories.UserRepositoryInterface) *ListingService {
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

	// Check if a vehicle with the same VIN already exists
	existingVehicle, err := s.vehicleRepo.GetByVIN(vehicle.VIN)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("error checking for existing vehicle: %w", err)
		}
		
		// Create vehicle
		if err := s.vehicleRepo.Create(vehicle); err != nil {
			return nil, fmt.Errorf("failed to create vehicle: %w", err)
		}
		// 'vehicle' now contains the ID of the newly created vehicle
	} else {
		// Vehicle with this VIN already exists, use its ID
		vehicle = existingVehicle // Use the existing vehicle's data (including its ID)
	}

	// 3. Link the listing to the (newly created or existing) vehicle
	listing.VehicleID = vehicle.ID

	// 4. Set initial status for the listing
	listing.Status = models.ListingStatusPending // Default to pending review

	// 5. Call the repository to persist the listing
	if err := s.repo.Create(listing); err != nil {
		// If listing creation fails, the vehicle might already be created,
		// but that's acceptable for this workflow. Rolling back the vehicle creation
		// would require a database transaction, which is more complex.
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// 6. Return the successfully created listing
	return listing, nil
}
