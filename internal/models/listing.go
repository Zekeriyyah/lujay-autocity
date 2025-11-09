package models

import (
	"time"

	"github.com/google/uuid"
)

type ListingStatus string

const (
	ListingStatusPending     ListingStatus = "pending_review"
	ListingStatusActive      ListingStatus = "active"
	ListingStatusRejected    ListingStatus = "rejected"
	ListingStatusSold        ListingStatus = "sold"
)

// Listing represents a specific listing of a vehicle on the platform.
type Listing struct {
	ID 			uuid.UUID	  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string        `json:"title" gorm:"size:255;not null"`
	Description string        `json:"description,omitempty" gorm:"type:text"`
	Price       float64       `json:"price" gorm:"not null;comment:Asking price in local currency"`
	Location    string        `json:"location" gorm:"size:255;not null"`
	Status      ListingStatus `json:"status" gorm:"default:pending_review;not null"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`

	// Foreign Keys
	SellerID  uuid.UUID `json:"seller_id" gorm:"type:uuid;not null"`
	VehicleID uuid.UUID `json:"vehicle_id" gorm:"type:uuid;not null"`

	// Relationships
	Seller    User     `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	Vehicle   Vehicle  `json:"vehicle" gorm:"foreignKey:VehicleID"`      
	Images    []Image  `json:"images,omitempty" gorm:"foreignKey:ListingID"`
}