package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshTokenString() (string, error) {

	tokenInt := make([]byte, 32)
	rand.Read(tokenInt)

	out := hex.EncodeToString(tokenInt)

	return out, nil
}