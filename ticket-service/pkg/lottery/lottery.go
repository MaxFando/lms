package lottery

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateTicketNumbers(count, max int) ([]string, error) {
	if count > max {
		return nil, fmt.Errorf("count can't be greater than max")
	}

	numSet := make(map[int]struct{})
	result := make([]string, 0, count)

	for len(result) < count {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			return nil, err
		}
		n := int(num.Int64()) + 1
		if _, exists := numSet[n]; !exists {
			numSet[n] = struct{}{}
			formatted := fmt.Sprintf("%02d", n)
			result = append(result, formatted)
		}
	}

	return result, nil
}
