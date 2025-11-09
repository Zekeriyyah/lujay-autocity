package repositories

import (
	"errors"
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/pkg"
	"gorm.io/gorm"
)

// UserRepositoryInterface defines the methods for interacting with User data.
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error) 
	GetByID(id string) (*models.User, error)     
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}


func (u *UserRepository) Create(user *models.User) error {

	if user == nil {
		return fmt.Errorf("user detail not provided")
	}

	if err := u.db.Create(user).Error; err != nil {

		// // Check for unique constraint violation
		// if erros.Is(err, gorm.ErrDuplicatedKey) {
		// 	return fmt.Errorf("user with email %s already exists", user.Email)
		// }
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}


func (u *UserRepository) GetByEmail(email string) (*models.User, error) {

	if u.db == nil {
        return nil, errors.New("database connection is not initialized")
    }

	user := models.User{}

	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}


func (u *UserRepository) GetByID(id string) (*models.User, error) {
	parsedID, err := pkg.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	user := models.User{}
	if err := u.db.First(&user, "id = ?", parsedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

