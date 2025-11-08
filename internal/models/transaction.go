package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
)

// Transaction represents a completed sale/purchase of a vehicle.
type Transaction struct {
	ID 				uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListingID       uuid.UUID `json:"listing_id" gorm:"type:uuid;uniqueIndex;not null"`
	BuyerID         uuid.UUID `json:"buyer_id" gorm:"type:uuid;not null"`
	SellerID        uuid.UUID `json:"seller_id" gorm:"type:uuid;not null"` // Denormalized for easy querying
	Amount          float64   `json:"amount" gorm:"not null;comment:Final sale price"`
	PaymentMethod   string    `json:"payment_method,omitempty" gorm:"size:100"`
	Status          TransactionStatus `json:"status" gorm:"default:pending;not null"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Listing Listing `json:"listing" gorm:"foreignKey:ListingID"`
	Buyer   User    `json:"buyer" gorm:"foreignKey:BuyerID"`
	Seller  User    `json:"seller" gorm:"foreignKey:SellerID"` // Denormalized for easy querying
}