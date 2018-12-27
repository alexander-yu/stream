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
			Vars:   stream.IntPtr(2),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})

	t.Run("fail: config with less than 2 vars is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Window: stream.IntPtr(3),
			Vars:   stream.IntPtr(1),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has less than 2 vars: %d < 2", 1))
	})

	t.Run("fail: Tuple with a negative exponent is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{{0, -1, 3, 4}},
			Window: stream.IntPtr(3),
			Vars:   stream.IntPtr(4),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a Tuple with a negative exponent of %d", -1))
	})

	t.Run("fail: Tuple with length != Vars is invalid", func(t *testing.T) {
		tuple := Tuple{0, 2, 3}
		config := &CoreConfig{
			Sums:   SumsConfig{tuple},
			Window: stream.IntPtr(3),
			Vars:   stream.IntPtr(4),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf(
			"config has a Tuple (%v) with length %d but Vars = %d",
			tuple,
			len(tuple),
			4,
		))
	})

	t.Run("fail: Tuple with all 0s is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{{0, 0, 0}},
			Window: stream.IntPtr(3),
			Vars:   stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config has a Tuple that is all 0s (i.e. skips all variables)")
	})

	t.Run("fail: config without Window is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{{1, 1, 3}},
			Vars: stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Window is not set")
	})

	t.Run("fail: config without Vars is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{{1, 1, 3}},
			Window: stream.IntPtr(3),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Vars is not set")
	})

	t.Run("pass: valid config is valid", func(t *testing.T) {
		config := &CoreConfig{
			Sums:   SumsConfig{{1, 1, 3}},
			Window: stream.IntPtr(3),
			Vars:   stream.IntPtr(3),
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
			Sums:   SumsConfig{{1, 1, 3}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{{1, 1, 3}},
			Vars:   stream.IntPtr(3),
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
			Sums:   SumsConfig{{3, 0, 0}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}
		mergedConfig, err := MergeConfigs(config)
		require.NoError(t, err)

		assert.Equal(t, config, mergedConfig)
	})

	t.Run("pass: multiple configs passed returns merged config", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{{1, 2, 3}, {2, 0, 0}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{{0, 2, 0}, {1, 2, 3}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}
		config3 := &CoreConfig{}
		config4 := &CoreConfig{
			Sums:   SumsConfig{{1, 2, 2}, {1, 1, 1}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}

		mergedConfig, err := MergeConfigs(config1, config2, config3, config4)
		require.NoError(t, err)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{{1, 2, 3}, {2, 0, 0}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}

		assert.Equal(t, expectedConfig, mergedConfig)
	})

	t.Run("fail: multiple configs passed fails if windows are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{{1, 2}, {2, 1}},
			Vars:   stream.IntPtr(2),
			Window: stream.IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{{0, 1}, {1, 1}},
			Vars:   stream.IntPtr(2),
			Window: stream.IntPtr(2),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing windows")
	})

	t.Run("fail: multiple configs passed fails if vars are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums:   SumsConfig{{1, 2}, {2, 1}},
			Vars:   stream.IntPtr(2),
			Window: stream.IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums:   SumsConfig{{0, 1, 1}, {1, 1, 2}},
			Vars:   stream.IntPtr(3),
			Window: stream.IntPtr(3),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing vars")
	})
}
