package config

import (
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config stores configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Version    string `mapstructure:"version"`
	DebugMode  bool   `mapstructure:"debug_mode"`
	HTTPServer struct {
		IP   string `mapstructure:"ip"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"http_server"`
	Token struct {
		SymmetricKey    string        `mapstructure:"symmetric_key"`
		AccessDuration  time.Duration `mapstructure:"access_duration"`
		RefreshDuration time.Duration `mapstructure:"refresh_duration"`
	} `mapstructure:"token"`
	MySQL struct {
		Log                 bool          `mapstructure:"log"`
		MaxOpenConnections  int           `mapstructure:"max_open_connections"`
		MaxIdleConnections  int           `mapstructure:"max_idle_connections"`
		MaxLifetimeDuration time.Duration `mapstructure:"max_lifetime_duration"`
		Read                struct {
			Addr     string `mapstructure:"addr"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Name     string `mapstructure:"name"`
		} `mapstructure:"read"`
		Write struct {
			Addr     string `mapstructure:"addr"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Name     string `mapstructure:"name"`
		} `mapstructure:"write"`
	} `mapstructure:"mysql"`
}

var (
	config Config
	once   sync.Once

	Set         = viper.Set
	Get         = viper.Get
	GetString   = viper.GetString
	GetInt      = viper.GetInt
	GetBool     = viper.GetBool
	GetFloat64  = viper.GetFloat64
	GetDuration = viper.GetDuration
)

// Init reads configuration from file or environment variables.
func Init() {
	InitWithPathAndType("config", "toml")
}

func InitWithPathAndType(configPath, configType string) {
	once.Do(func() {
		viper.AddConfigPath(configPath)
		viper.SetConfigName(GetEnv().String())
		viper.SetConfigType(configType)

		// read environment variables and automatically override values that it has read from configure file
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("cannot read config: ", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalln("failed to unmarshal from config: ", err)
		}

		viper.WatchConfig()

		viper.OnConfigChange(func(e fsnotify.Event) {
			if err := viper.Unmarshal(&config); err != nil {
				log.Fatalln("failed to unmarshal from config: ", err)
			}
		})
	})
}

// IsDebugMode returns whether the application is in debug mode.
func IsDebugMode() bool {
	return config.DebugMode
}

// GetVersion returns the version of the application.
func GetVersion() string {
	return config.Version
}

// GetConfig returns the configuration of the application.
func GetConfig() *Config {
	return &config
}

// GetIP returns the listen IP of the HTTP server.
func GetIP() string {
	return config.HTTPServer.IP
}

// GetIPWithDefault returns the listen IP of the HTTP server. If the IP is not set, it returns the default IP.
func GetIPWithDefault(defaultIP string) string {
	ip := GetIP()
	if ip == "" {
		return defaultIP
	}

	return ip
}

// GetPort returns the listen port of the HTTP server.
func GetPort() int {
	return config.HTTPServer.Port
}

// GetPortWithDefault returns the port of the HTTP server. If the port is not set, it returns the default port.
func GetPortWithDefault(defaultPort int) int {
	port := GetPort()
	if port == 0 {
		return defaultPort
	}

	return port
}

// GetTokenSymmetricKey returns the symmetric key of the token.
func GetTokenSymmetricKey() string {
	return config.Token.SymmetricKey
}

// GetTokenAccessDuration returns the duration of the access token.
func GetTokenAccessDuration() time.Duration {
	return config.Token.AccessDuration
}

// GetTokenRefreshDuration returns the duration of the refresh token.
func GetTokenRefreshDuration() time.Duration {
	return config.Token.RefreshDuration
}
