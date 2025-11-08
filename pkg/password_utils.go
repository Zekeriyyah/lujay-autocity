package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// CheckPassword verfies a password against it's hashed counterpart
func CheckPassword(hashPasswd, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd))
	if err != nil {
		Error(err, "password verification failed")
	}
	return err == nil
}