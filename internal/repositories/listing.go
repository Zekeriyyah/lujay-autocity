package repositories

import (
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"gorm.io/gorm"
)

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
		return fmt.Errorf("failed to create listing: %w", err)
	}
	return nil
}

// GetByID: retrieve listing by ID

func (r *ListingRepository) GetByID(id )