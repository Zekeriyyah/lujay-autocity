package repositories

import (
	"errors"
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/pkg"
	"gorm.io/gorm"
)

type VehicleRepositoryInterface interface {
	Create(vehicle *models.Vehicle) error
	CreateWithTx(tx *gorm.DB, vehicle *models.Vehicle)
	GetByID(id string) (*models.Vehicle, error)
	Update(vehicle *models.Vehicle) error
	Delete(id string) error
	GetByVIN(vin string) (*models.Vehicle, error)
	GetByVINWithTx(tx *gorm.DB, vin string) (*models.Vehicle, error)
}


type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository{
	return &VehicleRepository{db:db}
}

func(v *VehicleRepository) Create(vehicle *models.Vehicle) error {
	return v.CreateWithTx(v.db, vehicle)
}

func(v *VehicleRepository) CreateWithTx(tx *gorm.DB, vehicle *models.Vehicle) error {
	if err := tx.Create(vehicle).Error; err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

func(v *VehicleRepository) GetByID(id string)  (*models.Vehicle, error) {
	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return &models.Vehicle{}, fmt.Errorf("failed to parse string id to uuid: %w", err)
	}

	vehicle := &models.Vehicle{}
	if err := v.db.First(vehicle, "id = ?", parsedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Vehicle{}, fmt.Errorf("vehicle with id %s not found: %w", id, err)
		}
		return &models.Vehicle{}, fmt.Errorf("failed to get vehicle with id %s", id)
	}

	return vehicle, nil 
}

func(v *VehicleRepository) Update(vehicle *models.Vehicle) error {
	result := v.db.Updates(vehicle)
	if result.Error != nil {
		return fmt.Errorf("failed to update vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vehicle with id %s not found", vehicle.ID.String()) // Use UUID's String() method
	}
	return nil
}

func(v *VehicleRepository) Delete(id string) error {
	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return err
	}

	result := v.db.Delete(&models.Vehicle{}, "id = ?", parsedID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vehicle with id %s not found", id)
	}
	return nil
}

func (v *VehicleRepository) GetByVIN(vin string) (*models.Vehicle, error) {
	return v.GetByVINWithTx(v.db, vin)
}

func (v *VehicleRepository) GetByVINWithTx(tx *gorm.DB, vin string) (*models.Vehicle, error) {
	
	vehicle := &models.Vehicle{}
	
	if err := tx.Where("vin = ?", vin).First(vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get vehicle by VIN: %w", err)
	}
	return vehicle, nil
}