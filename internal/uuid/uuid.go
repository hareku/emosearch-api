package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUUID generates a unique id.
func GenerateUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("uuid random generation error: %w", err)
	}

	return uuid.String(), nil
}
