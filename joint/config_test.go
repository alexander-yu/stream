package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

func TestValidateConfig(t *testing.T) {
	t.Run("fail: config with a negative window is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Window: stream.IntPtr(-1),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})

	t.Run("fail: Tuple with a negative exponent is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{Tuple{0, -1, 3, 4}},
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a Tuple with a negative exponent of %d", -1))
	})

	t.Run("fail: Tuple with all 0s is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{Tuple{0, 0, 0}},
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config has a Tuple that is all 0s (i.e. skips all variables)")
	})

	t.Run("pass: empty config is valid", func(t *testing.T) {
		config := &CoreConfig{}
		err := validateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("pass: config with positive window is valid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 1, 3}},
			Window: stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.NoError(t, err)
	})
}

func TestSetConfigDefaults(t *testing.T) {
	t.Run("pass: empty config is overridden", func(t *testing.T) {
		config := &CoreConfig{}
		config = setConfigDefaults(config)

		expectedConfig := defaultConfig

		assert.Equal(t, expectedConfig, config)
	})

	t.Run("pass: provided fields are kept", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 1, 3}},
			Window: stream.IntPtr(3),
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 1, 3}},
			Window: stream.IntPtr(3),
		}

		assert.Equal(t, expectedConfig, config)
	})
}

func TestMergeConfigs(t *testing.T) {
	t.Run("fail: no configs passed is invalid", func(t *testing.T) {
		_, err := MergeConfigs()
		assert.EqualError(t, err, "no configs available to merge")
	})

	t.Run("pass: single config passed returns itself", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{Tuple{3, 0, 0}},
			Window: stream.IntPtr(3),
		}
		mergedConfig, err := MergeConfigs(config)
		require.NoError(t, err)

		assert.Equal(t, config, mergedConfig)
	})

	t.Run("pass: multiple configs passed returns union of sums and windows if all are compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 2, 3}, Tuple{2, 0, 0}},
			Window: stream.IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{Tuple{0, 2, 0}, Tuple{1, 2, 3}},
			Window: stream.IntPtr(3),
		}
		config3 := &CoreConfig{}

		mergedConfig, err := MergeConfigs(config1, config2, config3)
		require.NoError(t, err)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 2, 3}, Tuple{2, 0, 0}, Tuple{0, 2, 0}, Tuple{1, 2, 3}},
			Window: stream.IntPtr(3),
		}

		assert.Equal(t, expectedConfig, mergedConfig)
	})

	t.Run("fail: multiple configs passed fails if windows are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{Tuple{1, 2}, Tuple{2, 1}},
			Window: stream.IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{Tuple{0, 1}, Tuple{1, 1}},
			Window: stream.IntPtr(2),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing windows")
	})
}