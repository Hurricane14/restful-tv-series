package domain

import "github.com/google/uuid"

func NewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(UUID string) bool {
	_, err := uuid.Parse(UUID)
	return err == nil
}
