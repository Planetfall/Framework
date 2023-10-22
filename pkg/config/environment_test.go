package config_test

import (
	"testing"

	"github.com/planetfall/framework/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestOnCloud_withDev(t *testing.T) {
	// given
	envGiven := config.Development

	// when
	onCloudActual := envGiven.OnCloud()

	// then
	assert.Equal(t, false, onCloudActual)
}

func TestOnCloud_withPrd(t *testing.T) {
	// given
	envGiven := config.Production

	// when
	onCloudActual := envGiven.OnCloud()

	// then
	assert.Equal(t, true, onCloudActual)
}

func TestString(t *testing.T) {
	for envGiven := range config.EnvironmentMapping {
		envStringActual := envGiven.String()
		assert.IsType(t, "string", envStringActual)
	}
}
