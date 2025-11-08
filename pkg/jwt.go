package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID  uint   `json:"user_id"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func GeneratJWT(userID uint, t time.Time) (string, error) {
	claims := Claims{
		UserID:  userID,
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

	// Check for parsing errors
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
		if tokenHeader == "" {
			return "", fmt.Errorf("missing token")
		}

		if len(tokenSlice) != 2 || tokenSlice[1] == "" {
			return "", fmt.Errorf("invalid token")
		}

		token := tokenSlice[1]
		return token, nil
}