package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zekeriyyah/lujay-autocity/pkg"
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
		return fmt.Errorf("error checking for existing email: %w", err)
	}

	
	// Check if role is NOT admin AND NOT seller - set to default buyer if true
	if userData.Role == "" {
		userData.Role = types.RoleBuyer
	} else if userData.Role != types.RoleAdmin && userData.Role != types.RoleSeller && userData.Role != types.RoleBuyer {
		return fmt.Errorf("invalid role type specified")
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


// Login validates the user's credentials and returns a JWT token if successful.
func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	// 1. Get user by email from the database
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// 2. Compare the provided password with the hashed password from the database
	if  !user.VerifyPassword(password) {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// 3. Generate JWT token upon successful authentication
	expirationTime := time.Now().Add(24 * time.Hour)
	
	tokenString, err := pkg.GeneratJWT(user.ID, user.Role, expirationTime)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, user, nil
}
