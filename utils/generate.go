package utils

import (
	"crypto/rand"
	"math/big"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func GenerateShortLink(length int) (string, error) {
	res := make([]byte, length)
	aplahabetLen := big.NewInt(int64(len(base58Alphabet)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, aplahabetLen)
		if err != nil {
			return "", err
		}

		res[i] = base58Alphabet[num.Int64()]
	}

	return string(res), nil
}
