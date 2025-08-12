package util

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// CastID represents a cast identifier that encodes username, date, and unique ID
type CastID struct {
	Username string
	Date     string
	UniqueID string
}

// generateShortRandomID generates a short random ID using lo.RandomString
func generateShortRandomID() string {
	// Generate an 8-character random string with alphanumeric characters
	return lo.RandomString(8, lo.LowerCaseLettersCharset)
}

// GenerateCastID creates a new cast ID with the given username and date
func GenerateCastID(username, date string) (*CastID, error) {
	uniqueID := generateShortRandomID()

	return &CastID{
		Username: username,
		Date:     date,
		UniqueID: uniqueID,
	}, nil
}

// Encode converts the CastID to a short base64-encoded string
func (c *CastID) Encode() string {
	// Format: username_length:username:date:uniqueID
	data := fmt.Sprintf("%d:%s:%s:%s", len(c.Username), c.Username, c.Date, c.UniqueID)
	return base64.RawURLEncoding.EncodeToString([]byte(data))
}

// DecodeCastID decodes a short ID string back to CastID
func DecodeCastID(encoded string) (*CastID, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encoding: %w", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid format: expected 4 parts, got %d", len(parts))
	}

	return &CastID{
		Username: parts[1],
		Date:     parts[2],
		UniqueID: parts[3],
	}, nil
}

// FilePath returns the file path for this cast ID
func (c *CastID) FilePath() string {
	return fmt.Sprintf("%s/%s/%s.cast", c.Username, c.Date, c.UniqueID)
}

// String returns the encoded string representation
func (c *CastID) String() string {
	return c.Encode()
}
