package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ContainsError asserts that an error contains a string within its message.
func ContainsError(t *testing.T, err error, errString string) {
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), errString)
	}
}
