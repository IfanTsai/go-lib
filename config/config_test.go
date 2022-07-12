package config_test

import (
	"testing"

	"github.com/IfanTsai/go-lib/config"
	"github.com/stretchr/testify/require"
)

func TestGetPortWithDefault(t *testing.T) {
	require.Equal(t, 8888, config.GetPortWithDefault(8888))
}

func TestGetPort(t *testing.T) {
	config.SetEnv(config.TypeEnvDev)
	config.Init()
	require.Equal(t, 8080, config.GetPort())
}
