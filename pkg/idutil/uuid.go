package idutil

import (
	"crypto/rand"
	"fmt"
)

func UUID() string {
	// UUID is 16 bytes (128 bits)
	uuid := make([]byte, 16)

	// Fill with random bytes
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err) // Should never fail
	}

	// Set version (4) and variant bits according to RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	// Format as: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],   // 8 hex chars
		uuid[4:6],   // 4 hex chars
		uuid[6:8],   // 4 hex chars
		uuid[8:10],  // 4 hex chars
		uuid[10:16], // 12 hex chars
	)
}
