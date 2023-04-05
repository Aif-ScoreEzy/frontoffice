package helper

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAPIKey() string {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return err.Error()
	}

	apiKey := hex.EncodeToString(seed)

	return apiKey
}
