package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/microcosm-cc/bluemonday"
)

func GenerateHash() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Sanitize(str string) string {
	sanitizer := bluemonday.UGCPolicy()
	return sanitizer.Sanitize(str)
}