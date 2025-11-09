package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
)


func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	userID, ok := userIDVal.(uuid.UUID)
	return userID, ok
}

func GetUserRoleFromContext(c *gin.Context) (types.Role, bool) {
	roleVal, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := roleVal.(types.Role)
	return role, ok
}