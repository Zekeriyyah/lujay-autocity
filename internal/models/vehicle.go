package models

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle represents the core vehicle data.
type Vehicle struct {
	ID 			uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	VIN         string    `json:"vin" gorm:"uniqueIndex;size:17"`
	Make        string    `json:"make" gorm:"size:100;not null"`
	Model       string    `json:"model" gorm:"size:100;not null"`
	Year        int       `json:"year" gorm:"not null"`
	Mileage     uint      `json:"mileage" gorm:"comment:Odometer reading in miles/km"`
	EngineSize  string    `json:"engine_size,omitempty" gorm:"size:20;comment:Engine size (e.g., 2.0L)"`
	FuelType    string    `json:"fuel_type,omitempty" gorm:"size:50;comment:Gasoline, Diesel, Electric, Hybrid"`
	Transmission string   `json:"transmission,omitempty" gorm:"size:50;comment:Automatic, Manual"`
	BodyType    string    `json:"body_type,omitempty" gorm:"size:50;comment:Sedan, SUV, Truck, Coupe, etc."`
	Color       string    `json:"color,omitempty" gorm:"size:50"`
	Condition   string    `json:"condition,omitempty" gorm:"size:20;comment:New, Used, Certified Pre-Owned"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Listings []Listing `json:"-" gorm:"foreignKey:VehicleID"` // Omit from JSON
}