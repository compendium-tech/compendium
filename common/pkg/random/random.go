package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewRandomDigitCode(numberOfDigits uint8) (string, error) {
	upperBound := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(numberOfDigits)), nil)

	n, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		return "", err
	}

	format := fmt.Sprintf("%%0%dd", numberOfDigits)
	return fmt.Sprintf(format, n), nil
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
