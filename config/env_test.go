package config_test

import (
	"os"
	"testing"

	"github.com/IfanTsai/go-lib/config"
	"github.com/stretchr/testify/require"
)

func TestSetEnv(t *testing.T) {
	config.SetEnv(config.TypeEnvDev)
	require.Equal(t, config.TypeEnvDev, config.GetEnv())
	require.Equal(t, config.TypeEnvDev.String(), os.Getenv("ENV"))
}

func TestEnv(t *testing.T) {
	config.SetEnv(config.TypeEnvDev)
	require.Equal(t, config.TypeEnvDev, config.GetEnv())
	require.Equal(t, os.Getenv("ENV"), config.GetEnv().String())
}
