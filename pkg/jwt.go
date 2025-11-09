package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zekeriyyah/lujay-autocity/pkg/types"
)

var secret_key = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID  uuid.UUID   `json:"user_id"`
	Role	types.Role	`json:"role"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func GeneratJWT(userID uuid.UUID, role types.Role, t time.Time) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:	 role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret_key)
}

// validate if token is valid and retrieve the claims
func ValidateJWT(tokenStr string) (*Claims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		Info(tokenStr)
		Error(err, "error parsing token")
		return nil, err
	}

	// Check if token is valid and extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		Info("error invalid token or claims")
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil

}


func ExtractTokenStr(tokenHeader string) (string, error) {

	tokenHeader = strings.TrimSpace(tokenHeader)
	tokenSlice := strings.SplitN(tokenHeader, " ", 2)
		
	if len(tokenSlice) != 2 || tokenSlice[1] == "" {
		return "", fmt.Errorf("invalid token: bearer-token required")
	}

	token := tokenSlice[1]
	return token, nil
}