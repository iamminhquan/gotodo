// Package config manages application configuration using Viper and configdir.
package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

const (
	appName = "gotodo"

	// Viper config keys.
	KeyStorageType     = "storage_type"
	KeyDataDir         = "data_dir"
	KeyDefaultPriority = "default_priority"
	KeyUseColor        = "use_color"
	KeyTasksFileName   = "tasks_file_name"
)

// Config holds strongly-typed application settings loaded from the config file
// and / or environment variables.
type Config struct {
	StorageType     string `mapstructure:"storage_type"`
	DataDir         string `mapstructure:"data_dir"`
	DefaultPriority string `mapstructure:"default_priority"`
	UseColor        bool   `mapstructure:"use_color"`
	TasksFileName   string `mapstructure:"tasks_file_name"`
}

// appConfigDir returns the OS-appropriate configuration directory for gotodo.
func appConfigDir() string {
	return configdir.LocalConfig(appName)
}

// appDataDir returns the OS-appropriate data directory for gotodo.
func appDataDir() string {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("APPDATA")
		if base == "" {
			base, _ = os.UserHomeDir()
		}

		return filepath.Join(base, appName, "data")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", appName)
	default: // linux and others
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg != "" {
			return filepath.Join(xdg, appName)
		}

		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "share", appName)
	}
}

// Init loads configuration from the config file and environment.
// Must be called once at application start-up (e.g. in PersistentPreRunE).
func Init() (*Config, error) {
	viper.SetDefault(KeyStorageType, "json")
	viper.SetDefault(KeyDataDir, appDataDir())
	viper.SetDefault(KeyDefaultPriority, "medium")
	viper.SetDefault(KeyUseColor, true)
	viper.SetDefault(KeyTasksFileName, "tasks.json")

	// Allow environment variable overrides: GOTODO_DATA_DIR etc.
	viper.SetEnvPrefix("GOTODO")
	viper.AutomaticEnv()

	configDir := appConfigDir()
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		// Non-fatal: continue with defaults if we can't create the config dir.
		_ = err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Ignore "file not found" – the app works fine on first run without a
	// config file.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfigDir returns the directory where the config file is stored.
func GetConfigDir() string {
	return appConfigDir()
}

// GetTasksFilePath returns the full path to the tasks JSON file.
func GetTasksFilePath(cfg *Config) string {
	return filepath.Join(cfg.DataDir, cfg.TasksFileName)
}

// EnsureDataDir creates the data directory if it does not already exist.
func EnsureDataDir(cfg *Config) error {
	return os.MkdirAll(cfg.DataDir, 0o755)
}
