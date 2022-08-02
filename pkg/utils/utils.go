package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64EncodeInt(n int) string {
	return Base64EncodeBytes([]byte(strconv.Itoa(n)))
}

func Sha1Hash(s string) string {
	hash := sha1.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
