package utils

import (
	"crypto/rand"
	"encoding/base64"
)

type CtxKey string

const CtxURLPathValueKey = CtxKey("url_value")
func GenerateNonce() string {
	nonceBytes := make([]byte, 16)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		// TODO: Handle the error appropriately
		return ""
	}
	return base64.StdEncoding.EncodeToString(nonceBytes)
}
