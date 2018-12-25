package mathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorial(t *testing.T) {
	assert.Equal(t, 3628800, factorial(10))
}

func TestSign(t *testing.T) {
	assert.Equal(t, 1, Sign(0))
	assert.Equal(t, -1, Sign(-1))
	assert.Equal(t, -1, Sign(1))
}

func TestBinom(t *testing.T) {
	assert.Equal(t, 10, Binom(5, 2))
	assert.Equal(t, 20, Binom(20, 1))
	assert.Equal(t, 1, Binom(1500, 0))
}
