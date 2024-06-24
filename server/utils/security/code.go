package security

import (
	"fmt"
	"math/rand"
)

// GenerateCode generates a random user code
func GenerateCode() string {
	randNum := rand.Intn(10000)
	code := fmt.Sprintf("%04d", randNum)
	return code
}
