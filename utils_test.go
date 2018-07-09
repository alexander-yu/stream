package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorial(t *testing.T) {
	assert.Equal(t, 3628800, factorial(10))
}

func TestSign(t *testing.T) {
	assert.Equal(t, 1, sign(0))
	assert.Equal(t, -1, sign(-1))
	assert.Equal(t, -1, sign(1))
}

func TestBinom(t *testing.T) {
	assert.Equal(t, 10, binom(5, 2))
	assert.Equal(t, 20, binom(20, 1))
	assert.Equal(t, 1, binom(1500, 0))
}
