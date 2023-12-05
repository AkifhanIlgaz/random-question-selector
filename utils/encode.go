package utils

import (
	"encoding/base64"
	"fmt"
)

func Encode(s string) string {
	return string(base64.StdEncoding.EncodeToString([]byte(s)))
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base 64 decode: %w", err)
	}

	return string(data), nil
}
