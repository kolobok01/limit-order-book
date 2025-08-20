package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func Sign(timestamp, msgType, seqNum, key, targetCompId, passphrase, secret string) string {
	// Create pre-hash string: "<timestamp>A1<key><targetCompId><passphrase>"
	prehash := fmt.Sprintf("%s%s%s%s%s%s", timestamp, msgType, seqNum, key, targetCompId, passphrase)

	// Decode base64 secret
	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}

	// Create HMAC
	h := hmac.New(sha256.New, secretBytes)
	h.Write([]byte(prehash))

	// Encode signature to base64
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.TrimRight(signature, "=")
}
