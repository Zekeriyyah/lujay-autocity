package repositories

import (
	"errors"
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/pkg"
	"gorm.io/gorm"
)

type ListingRepositoryInterface interface {
	Create(listing *models.Listing) error
	GetByID(id string) (*models.Listing, error) 
	Update(listing *models.Listing) error
	Delete(id string) error 
	GetBySellerID(sellerID string) ([]*models.Listing, error) 
	
	GetAll() ([]*models.Listing, error)
	GetByStatus(status models.ListingStatus) ([]*models.Listing, error)
}

type ListingRepository struct {
	db *gorm.DB
}

func NewListingRepository(db *gorm.DB) *ListingRepository {
	return &ListingRepository{
		db: db,
	}
}

// Create: inserts a listing into the db
func (r *ListingRepository) Create(listing *models.Listing) error {
	if err := r.db.Create(listing).Error; err != nil {
		return fmt.Errorf("❌ failed to create listing: %w", err)
	}
	return nil
}

// GetByID: retrieve listing by ID
func (r *ListingRepository) GetByID(id string) (*models.Listing, error) {

	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return &models.Listing{}, err
	}

	listing := &models.Listing{}
	
	// Preload Seller, Vehicle and Image
	if err := r.db.Preload("Seller").Preload("Vehicle").Preload("Images").First(listing, "id = ?", parsedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("❌ listing with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}
	return listing, nil
}

// Update: update fields in a listing
func (r *ListingRepository) Update(listing *models.Listing) error {
	result := r.db.Save(listing)
	if result.Error != nil {
		return fmt.Errorf("❌ failed to update listing: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("❌ listing with id %s not found", listing.ID.String()) 
	}
	
	return nil
}

func (r *ListingRepository) Delete(id string) error {

	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return err
	}

	result := r.db.Delete(&models.Listing{}, "id = ?", parsedID) 
	if result.Error != nil {
		return fmt.Errorf("❌ failed to delete listing: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("❌ listing with id %s not found", id)
	}
	
	return nil
}

func(r *ListingRepository) GetBySellerID(id string) ([]*models.Listing, error) {
	
	parseSellerID, err := pkg.StringToUUID(id) 
	if err != nil {
		return nil, fmt.Errorf("❌ failed to parse seller id to uuid: %w", err)
	}
	
	listings := []*models.Listing{}

	if err = r.db.Where("seller_id = ?", parseSellerID).Preload("Vehicle").Preload("Images").Find(&listings).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to get listings by seller ID: %w", err)
	}

	return listings, nil
}

func(r *ListingRepository) GetAll() ([]*models.Listing, error) {

	listings := []*models.Listing{}

	err := r.db.Preload("Seller").Preload("Vehicle").Preload("Images").Find(&listings).Error
	if err != nil {
		return nil, fmt.Errorf("❌ failed to get all listings: %w", err)
	}
	return listings, nil
}

func (r *ListingRepository) GetByStatus(status models.ListingStatus) ([]*models.Listing, error) {
	listings := []*models.Listing{}

	err := r.db.Where("status=?", status).Preload("Seller").Preload("Vehicle").Preload("Images").Find(&listings).Error
	if err != nil {
		return nil, fmt.Errorf("❌ failed to get %s listings: %w", status, err) 
	}
	return listings, nil
}