package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Manager handles configuration management with support for file and environment variable overrides
type Manager struct {
	viper *viper.Viper
	log   *logrus.Entry
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		viper: viper.New(),
		log: logrus.WithFields(map[string]interface{}{
			"module": "config",
		}),
	}
}

func (m *Manager) Viper() *viper.Viper {
	return m.viper
}

func (m *Manager) Set(key, value string) {
	m.viper.Set(key, value)
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(configFile string) error {
	if configFile == "" {
		m.log.Warn("No config file specified")
		return nil
	}

	m.viper.SetConfigFile(configFile)
	if err := m.viper.ReadInConfig(); err != nil {
		m.log.Warnf("Config file not found: %s, using environment variables", configFile)
		return err
	}

	m.log.Infof("Loaded config from file: %s", configFile)
	return nil
}

// SetupEnvironmentOverrides sets up environment variable overrides using Viper's built-in support
func (m *Manager) SetupEnvironmentOverrides() {
	// Enable automatic environment variable lookup
	m.viper.AutomaticEnv()

	// Replace dots with underscores for environment variable names
	// Example: server.port -> SERVER_PORT (when using prefix)
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// SetEnvPrefix sets a prefix for environment variables
// Example: SetEnvPrefix("APP") means APP_SERVER_PORT maps to server.port
func (m *Manager) SetEnvPrefix(prefix string) {
	m.viper.SetEnvPrefix(prefix)
	m.log.Infof("Environment variable prefix set to: %s", prefix)
}

// BindEnv binds environment variables to configuration keys
func (m *Manager) BindEnv(key, envVar string) {
	m.viper.BindEnv(key, envVar)
}

// BindEnvs binds multiple environment variables to configuration keys
func (m *Manager) BindEnvs(bindings map[string]string) {
	for key, envVar := range bindings {
		m.viper.BindEnv(key, envVar)
	}
}

// GetString returns a string configuration value
func (m *Manager) GetString(key string) string {
	return m.viper.GetString(key)
}

// GetInt returns an integer configuration value
func (m *Manager) GetInt(key string) int {
	return m.viper.GetInt(key)
}

// GetBool returns a boolean configuration value
func (m *Manager) GetBool(key string) bool {
	return m.viper.GetBool(key)
}

// GetStringSlice returns a string slice configuration value
func (m *Manager) GetStringSlice(key string) []string {
	return m.viper.GetStringSlice(key)
}

// UnmarshalKey unmarshals a configuration key into a struct
func (m *Manager) UnmarshalKey(key string, rawVal interface{}) error {
	return m.viper.UnmarshalKey(key, rawVal)
}

// UnmarshalKeyWithEnv unmarshals a configuration key into a struct
// and automatically syncs environment variable values before unmarshaling
// envMappings: map[configKey]envVar (e.g., map["server.port"]="SERVER_PORT")
func (m *Manager) UnmarshalKeyWithEnv(key string, rawVal interface{}, envMappings map[string]string) error {
	// Auto-sync environment variables directly
	for configKey, envVar := range envMappings {
		if envValue := os.Getenv(envVar); envValue != "" {
			m.viper.Set(configKey, envValue)
			m.log.Debugf("Synced env %s=%s to config %s", envVar, envValue, configKey)
		}
	}
	return m.viper.UnmarshalKey(key, rawVal)
}

// Unmarshal unmarshals the entire configuration into a struct
func (m *Manager) Unmarshal(rawVal interface{}) error {
	return m.viper.Unmarshal(rawVal)
}

// GetViper returns the underlying viper instance for advanced usage
func (m *Manager) GetViper() *viper.Viper {
	return m.viper
}

// LogConfigValue logs a configuration value for debugging
func (m *Manager) LogConfigValue(key string) {
	value := m.viper.GetString(key)
	envVar := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	envValue := os.Getenv(envVar)

	m.log.Infof("Config %s: %s (env %s: %s)", key, value, envVar, envValue)
}
