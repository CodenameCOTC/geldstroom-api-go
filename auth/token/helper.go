package token

import (
	"crypto/rand"
	"fmt"
)

func generateRandomToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
