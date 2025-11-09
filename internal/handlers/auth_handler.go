package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zekeriyyah/lujay-autocity/internal/models"
	"github.com/zekeriyyah/lujay-autocity/internal/services"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
)



type AuthHandler struct {
	service   *services.AuthService
	validator *validator.Validate
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator.New(), 
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email    string      `json:"email" validate:"required,email"`
		Name     string      `json:"name" validate:"required"`
		Phone    string      `json:"phone"`
		Password string      `json:"password" validate:"required,min=6"`
		Role     types.Role `json:"role"`
	}

	// Bind JSON request to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	// Map input to the User model
	user := &models.User{
		Email:    input.Email,
		Name:     input.Name,
		Phone:    input.Phone,
		Password: input.Password, 
		Role:     input.Role,   
	}

	if err := h.service.Register(user); err != nil {
		log.Printf("Error registering user: %v", err)
		
		if err.Error() == fmt.Sprintf("user with email %s already exists", input.Email) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// For other errors, return a generic message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed: "+ err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"body": gin.H{
			"email": user.Email,
			"name": user.Name,
			"phone": user.Phone,
			"role": user.Role,
		},
	})
}