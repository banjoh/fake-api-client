package accounts

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff(t *testing.T) {
	backoffs := []int{1, 2, 4, 8, 16, 16, 16, 16}

	curr := initialBackoffSeconds
	for _, v := range backoffs {
		next := nextBackoff(curr)
		curr = next

		assert.Equal(t, math.Floor(next), float64(v))
	}
}
