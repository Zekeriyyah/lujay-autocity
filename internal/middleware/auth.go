package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/pkg"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
)

func AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// access the authHeader for the token
	
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token string from header
		token, err := pkg.ExtractTokenStr(tokenHeader)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Extract claim from token
		claims, err := pkg.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		
		c.Next()
	})
}

// RBAC middleware checks if the authenticated user's role matches the required role
func RBAC(requiredRoles ...types.Role) gin.HandlerFunc{
	return gin.HandlerFunc(func(c *gin.Context) {
		userRoleVal, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user role not found in context"})
			c.Abort()
			return
		}

		// assert userRoleVal's type
		userRole, ok := userRoleVal.(types.Role)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user role type assertion failed"})
			c.Abort()
			return
		}


		for _, requiredRole := range requiredRoles {
			log.Println(requiredRole, userRoleVal)
			if userRole == requiredRole {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permission"})
		c.Abort()
	}) 
}


