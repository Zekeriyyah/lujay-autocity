package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zekeriyyah/lujay-autocity/pkg"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// access the authHeader for the token
	
		tokenHeader := c.GetHeader("Authorization")
		token, err := pkg.ExtractTokenStr(tokenHeader)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error parsing token header"})
			c.Abort()
			return
		}

		claims, err := pkg.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		
		c.Next()
	}
}
