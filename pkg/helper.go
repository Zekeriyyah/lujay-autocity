package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// ParseStringifiedJSON attempts to parse a string that represents JSON data
func ParseStringifiedJSON(jsonString string) (map[string]any, error) {
	if jsonString == "" {
		return map[string]any{}, nil
	}

	var parsedData map[string]any
	var unmarshalingError error

	if err := json.Unmarshal([]byte(jsonString), &parsedData); err == nil {
		return parsedData, nil
	} else {
		unmarshalingError = err
	}

	unquotedString, unquoteErr := strconv.Unquote(jsonString)
	if unquoteErr != nil {
		return nil, fmt.Errorf("failed to unquote and unmarshal JSON string. Unquote error: %w. Initial unmarshal error: %v", unquoteErr, unmarshalingError)
	}

	if err := json.Unmarshal([]byte(unquotedString), &parsedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal unquoted JSON string '%s': %w", unquotedString, err)
	}

	return parsedData, nil
}