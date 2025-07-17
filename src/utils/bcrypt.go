package utils

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func ParseUUID(uuidStr string) (uuid.UUID, error) {
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return u, nil
}
