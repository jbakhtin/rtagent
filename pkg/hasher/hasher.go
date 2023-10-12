package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func CalcHash(value string, salt []byte) (string, error) {
	h := hmac.New(sha256.New, []byte(salt))

	_, err := h.Write([]byte(value))
	if err != nil {
		return "", err
	}

	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst), nil
}