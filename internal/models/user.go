package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleSeller  Role = "seller"
	RoleBuyer   Role = "buyer"
)

// User represents the user profile in the system.
type User struct {
	ID 		  uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Role      Role      `json:"role" gorm:"default:buyer;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships (Optional, for GORM)
	Listings      []Listing      `json:"-" gorm:"foreignKey:SellerID"`
	Transactions  []Transaction  `json:"-" gorm:"foreignKey:BuyerID"`
}