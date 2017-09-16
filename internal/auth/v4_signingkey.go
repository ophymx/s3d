package auth

import "crypto/subtle"

type SigningKey []byte

func (key SigningKey) Sign(str string) string {
	return string(key.sign(str))
}

func (key SigningKey) Verify(str, signature string) bool {
	return subtle.ConstantTimeCompare(key.sign(str), []byte(signature)) == 1
}

func (key SigningKey) sign(str string) []byte {
	return hexEncode(hmacSha256([]byte(key), str))
}
