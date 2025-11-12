package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/pkg"
)

type InspectionStatus string

const (
	InspectionStatusPending   InspectionStatus = "pending"
	InspectionStatusApproved  InspectionStatus = "approved"
	InspectionStatusRejected  InspectionStatus = "rejected"
)

// Inspection represents the vetting report for a listed vehicle.
type Inspection struct {
	ID 				uuid.UUID		   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListingID       uuid.UUID          `json:"listing_id" gorm:"type:uuid;uniqueIndex;not null"`
	InspectorID     uuid.UUID          `json:"inspector_id" gorm:"type:uuid;not null"`
	InspectionDate  time.Time          `json:"inspection_date,omitempty"`
	ConditionRating int                `json:"condition_rating,omitempty" gorm:"comment:Overall condition score (e.g., 1-10)"`
	Findings        map[string]any     `json:"findings,omitempty" gorm:"type:jsonb;comment:Detailed inspection notes (e.g., dents, scratches, mechanical issues)"`
	ReportURL       string             `json:"report_url,omitempty" gorm:"size:500;comment:Link to the full inspection report PDF/image"`
	Status          InspectionStatus   `json:"status" gorm:"default:pending;not null"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`

	// Relationships
	Listing   Listing `json:"listing" gorm:"foreignKey:ListingID"`
	Inspector User    `json:"inspector,omitempty" gorm:"foreignKey:InspectorID"`
}


// Inspection represents the vetting report for a listed vehicle.
type InspectionFetchInput struct {
	ID              uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListingID       uuid.UUID          `json:"listing_id" gorm:"type:uuid;uniqueIndex;not null"`
	InspectorID     uuid.UUID          `json:"inspector_id" gorm:"type:uuid;not null"`
	InspectionDate  time.Time          `json:"inspection_date,omitempty"`
	ConditionRating int                `json:"condition_rating,omitempty" gorm:"comment:Overall condition score (e.g., 1-10)"`

	Findings        string             `json:"findings,omitempty" gorm:"comment:Detailed inspection notes (e.g., dents, scratches, mechanical issues)"` // GORM tag might not be needed for reading, but kept for clarity if used in other contexts

	ReportURL       string             `json:"report_url,omitempty" gorm:"size:500;comment:Link to the full inspection report PDF/image"`
	Status          InspectionStatus   `json:"status" gorm:"default:pending;not null"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`

	
	Listing   Listing `json:"listing" gorm:"foreignKey:ListingID"`
	Inspector User    `json:"inspector,omitempty" gorm:"foreignKey:InspectorID"`
}


func (i *InspectionFetchInput) ParseInspectionInputToModel() ( *Inspection, error) {
	findings, err := pkg.ParseStringifiedJSON(i.Findings)
	if err != nil {
		return nil, err
	}

	inspection := &Inspection{
		ID: i.ID,
		ListingID: i.ListingID,
		InspectorID: i.InspectorID,
		InspectionDate: i.InspectionDate,
		ConditionRating: i.ConditionRating,

		Findings: findings,
		ReportURL: i.ReportURL,
		Status: i.Status,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,		
		Listing: i.Listing,
		Inspector:	i.Inspector,
	}

	return inspection, nil
}