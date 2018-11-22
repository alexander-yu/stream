package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	t.Run("fail: config with a non-nil empty map is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{},
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config sums map is not nil but empty")
	})

	t.Run("fail: config with a nonpositive window is invalid", func(t *testing.T) {
		config := &CoreConfig{
			Window: IntPtr(0),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config window is nonpositive")
	})

	t.Run("pass: empty config is valid", func(t *testing.T) {
		config := &CoreConfig{}
		err := validateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("pass: config with non-empty map and positive window is valid", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{2: true},
			Window: IntPtr(3),
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

	t.Run("pass: empty sums is overridden, provided window is kept", func(t *testing.T) {
		config := &CoreConfig{
			Window: IntPtr(3),
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums:   defaultConfig.Sums,
			Window: IntPtr(3),
		}

		assert.Equal(t, expectedConfig, config)
	})

	t.Run("pass: empty window is overridden, provided sums is kept", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{3: true},
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums:   SumsConfig{3: true},
			Window: defaultConfig.Window,
		}

		assert.Equal(t, expectedConfig, config)
	})

	t.Run("pass: provided window and sums are kept", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{3: true},
			Window: IntPtr(3),
		}
		config = setConfigDefaults(config)

		expectedConfig := &CoreConfig{
			Sums: SumsConfig{3:, true},
			Window: IntPtr(3),
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
			Sums: SumsConfig{3: true},
			Window: IntPtr(3),
		}
		mergedConfig, err := MergeConfigs(config)
		require.NoError(t, err)

		assert.Equal(t, config, mergedConfig)
	})

	t.Run("pass: multiple configs passed returns union of sums and windows if all are compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums: SumsConfig{1: true, 2: true},
			Window: IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums: SumsConfig{2: true, 3: true},
			Window: IntPtr(3),
		}
		config3 := &CoreConfig{}

		mergedConfig, err := MergeConfigs(config1, config2)
		require.NoError(t, err)

		expectedConfig := &CoreConfig{
			Sums: SumsConfig{1: true, 2: true, 3: true},
			Window: IntPtr(3),
		}

		assert.Equal(t, expectedConfig, mergedConfig)
	})

	t.Run("fail: multiple configs passed fails if windows are not compatible", func(t *testing.T) {
		config1 := &CoreConfig{
			Sums: SumsConfig{1: true, 2: true},
			Window: IntPtr(3),
		}
		config2 := &CoreConfig{
			Sums: SumsConfig{2: true, 3: true},
			Window: IntPtr(2),
		}

		_, err := MergeConfigs(config1, config2)
		assert.EqualError(t, err, "configs have differing windows")
	})
}
