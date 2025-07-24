package random

import (
	"math"
	"math/big"
	"strconv"

	"crypto/rand"
)

func NewRandomDigitCode(numberOfDigits uint8) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(10, float64(numberOfDigits)))))
	if err != nil {
		return "", nil
	}

	return strconv.FormatInt(n.Int64(), 10), nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func NewRandomString(size uint) (string, error) {
	code := make([]byte, size)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		code[i] = charset[idx.Int64()]
	}

	return string(code), nil
}
