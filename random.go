package password_generator

import (
	"crypto/rand"
	"math/big"
)

// randBetween generates a random big int between min and max.
//
// Parameters:
//   - min: The minimum value of the range.
//   - max: The maximum value of the range.
//
// Returns:
//   - A random int64 value between min and max.
//   - An error if the random number generation fails.
func randBetween(min, max int) (int64, error) {
	diff := big.NewInt(int64(max - min))
	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, err
	}

	// Add min to shift the range
	n.Add(n, big.NewInt(int64(min)))
	return n.Int64(), nil
}

// randomNumber generates a random number within the given max limit.
//
// Parameters:
//   - maxLimit: The maximum value of the generated number, not inclusive.
//
// Returns:
//   - int64: The generated number.
//   - error: An error if the operation fails.
func randomNumber(maxLimit int) (int64, error) {
	max := big.NewInt(int64(maxLimit))
	index, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return index.Int64(), nil
}
