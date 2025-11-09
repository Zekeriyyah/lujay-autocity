package services

import (
	"errors"
	"fmt"

	"github.com/zekeriyyah/lujay-autocity/pkg/types"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/repositories"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo repositories.UserRepositoryInterface
	jwtSecret string
}


func NewAuthService(repo repositories.UserRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: repo,
		jwtSecret: jwtSecret,
	}
}

func (a *AuthService) Register(userData *models.User) error {
	// check if email already exists
	_, err := a.userRepo.GetByEmail(userData.Email)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", userData.Email)
	}
	
	// for errors other than RecordNotFound
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking for existing email")
	}
	
// Check if role is NOT admin AND NOT seller - set to default buyer if true
if userData.Role != types.RoleAdmin && userData.Role != types.RoleSeller {
    userData.Role = types.RoleBuyer
}

	// hash password and set it
	err = userData.SetPassword(userData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	} 

	// save the user
	if err = a.userRepo.Create(userData); err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
} 