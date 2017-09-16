package auth

import (
	"crypto/subtle"
	"encoding/base64"
)

type SigningKeyV2 []byte

func (key SigningKeyV2) Sign(str string) string {
	return string(key.sign(str))
}

func (key SigningKeyV2) Verify(str, signature string) bool {
	return subtle.ConstantTimeCompare(key.sign(str), []byte(signature)) == 1
}

func (key SigningKeyV2) sign(str string) []byte {
	return []byte(base64.StdEncoding.EncodeToString(hmacSha1([]byte(key), str)))
}
