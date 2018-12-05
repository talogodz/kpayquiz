package rand

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func RandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
		return ""
	}
	return hex.EncodeToString(bytes)
}
