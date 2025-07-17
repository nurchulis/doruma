package utils

import (
	"fmt"

	"time"

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

// ParseDatetime parses a datetime string in RFC3339 format and returns a time.Time object.
func ParseDatetime(datetimeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, datetimeStr)
}
