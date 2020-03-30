package passwd

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomBase64String(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(b)[:length], nil
}
