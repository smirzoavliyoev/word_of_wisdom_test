package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)

// randomBytes reads n cryptographically secure pseudo-random numbers.
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// base64EncodeBytes
func Base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// base64EncodeInt
func Base64EncodeInt(n int) string {
	return Base64EncodeBytes([]byte(strconv.Itoa(n)))
}

// sha1Hash
func Sha1Hash(s string) string {
	hash := sha1.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
