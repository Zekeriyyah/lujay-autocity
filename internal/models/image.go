package models

import (
	"time"

	"github.com/google/uuid"
)

// Image represents a file associated with a listing.
type Image struct {
	ID 		  uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListingID uuid.UUID `json:"listing_id" gorm:"type:uuid;not null"` 
	URL       string    `json:"url" gorm:"size:500;not null"`
	Order     int       `json:"order,omitempty" gorm:"default:0;comment:Display order for the image gallery"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Listing Listing `json:"-" gorm:"foreignKey:ListingID"`
}