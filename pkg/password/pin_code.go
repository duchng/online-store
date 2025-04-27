package password

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateRandomPinCode(pinLength int) (string, error) {
	const digits = "0123456789"
	pin := make([]byte, pinLength)
	for i := 0; i < pinLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return strings.Repeat("0", pinLength), err
		}
		pin[i] = digits[num.Int64()]
	}
	return string(pin), nil
}
