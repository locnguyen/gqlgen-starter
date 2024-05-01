package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationConfig_IsDevelopment(t *testing.T) {
	devCfg := ApplicationConfig{
		GoEnv: "development",
	}

	prodCfg := ApplicationConfig{
		GoEnv: "production",
	}

	assert.True(t, devCfg.IsDevelopment(), "true when GoEnv=development")
	assert.False(t, prodCfg.IsDevelopment(), "false GoEnv=production")
}

func TestApplicationConfig_IsProduction(t *testing.T) {
	devCfg := ApplicationConfig{
		GoEnv: "development",
	}

	prodCfg := ApplicationConfig{
		GoEnv: "production",
	}

	assert.False(t, devCfg.IsProduction(), "true when GoEnv=development")
	assert.True(t, prodCfg.IsProduction(), "false GoEnv=production")
}
