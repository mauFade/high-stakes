package util

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// GenerateKSUID generates a KSUID-like identifier
// Simplified version: timestamp (base64) + random payload (base64)
func GenerateKSUID() (string, error) {
	// Get current timestamp (nanoseconds for uniqueness)
	timestamp := time.Now().UnixNano()

	// Generate random payload (16 bytes)
	payload := make([]byte, 16)
	if _, err := rand.Read(payload); err != nil {
		return "", err
	}

	// Combine timestamp (8 bytes) + payload (16 bytes) = 24 bytes
	data := make([]byte, 24)
	// Encode timestamp as 8 bytes (big endian)
	for i := 7; i >= 0; i-- {
		data[i] = byte(timestamp & 0xff)
		timestamp >>= 8
	}
	copy(data[8:], payload)

	// Base64 URL encode (removes padding, uses URL-safe characters)
	ksuid := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
	
	// Ensure it's exactly 32 characters (or close to it)
	// KSUID format is typically 27 chars, but we'll use base64 for simplicity
	if len(ksuid) > 27 {
		ksuid = ksuid[:27]
	}
	
	return ksuid, nil
}

