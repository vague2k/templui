package utils

import (
	"crypto/rand"
	"encoding/base64"
)

type CtxKey string

const CtxURLPathValueKey = CtxKey("url_value")

func GenerateNonce() string {
	try := 0
Retry:
	try++

	nonceBytes := make([]byte, 16)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		if try < 2 {
			goto Retry
		}

		return ""
	}
	return base64.StdEncoding.EncodeToString(nonceBytes)
}
