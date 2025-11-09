package pkg

import (
	"fmt"

	"github.com/google/uuid"
)

// StringToUUID: Helper to convert string ID to UUID
func StringToUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}
	return parsedUUID, nil
}
