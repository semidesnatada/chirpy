package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	hash, hErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	if hErr != nil {
		return "", hErr
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {

	success := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return success
}