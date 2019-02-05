package moment

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

	t.Run("fail: config with a decay not in (0, 1) is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a decay of %f, which is not in (0, 1)", 0.))

		config = &CoreConfig{
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(1),
		}
		err = validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a decay of %f, which is not in (0, 1)", 1.))
	})

	t.Run("fail: config with a set decay and nonzero window is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config cannot have Decay set with a nonzero window")
	})

	t.Run("fail: config with a nonpositive central moment is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   map[int]bool{-1: true},
			Window: stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a nonpositive central moment of %d", -1))
	})

	t.Run("fail: config without Window is invalid", func(t *testing.T) {
		config := &CoreConfig{}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Window is not set")
	})

	t.Run("pass: valid config is valid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{2: true},
			Window: stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.NoError(t, err)

		config = &CoreConfig{
			Sums:   SumsConfig{2: true},
			Window: stream.IntPtr(0),
			Decay:  stream.FloatPtr(0.3),
		}
		err = validateConfig(config)
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
			Sums:   SumsConfig{3: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{3: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
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
			Sums:   SumsConfig{3: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		mergedConfig, err := MergeConfigs(config)
		require.NoError(t, err)

		assert.Equal(t, config, mergedConfig)
	})

	t.Run("pass: multiple configs passed returns union of sums and windows if all are compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{1: true, 2: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{2: true, 3: true},
			Window: stream.IntPtr(3),
		}
		config3 := &CoreConfig{}

		mergedConfig, err := MergeConfigs(config1, config2, config3)
		require.NoError(t, err)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{1: true, 2: true, 3: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}

		assert.Equal(t, expectedConfig, mergedConfig)
	})

	t.Run("fail: multiple configs passed fails if windows are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{1: true, 2: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{2: true, 3: true},
			Window: stream.IntPtr(2),
			Decay:  stream.FloatPtr(0.3),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing windows")
	})

	t.Run("fail: multiple configs passed fails if decays are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{1: true, 2: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{2: true, 3: true},
			Window: stream.IntPtr(3),
			Decay:  stream.FloatPtr(0.5),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing decays")
	})
}
