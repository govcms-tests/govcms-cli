package settings

import (
	"os/user"
	"path/filepath"

	"github.com/govcms-tests/govcms-cli/pkg/utils"
	"github.com/spf13/viper"
)

// Config is a struct of configurable options
type Config struct {
	Workspace string            `yaml:"workspace"` // The project workspace
	OS        string            `yaml:"os"`        // Operating system information
	Tokens    map[string]string `yaml:"tokens"`    // Access tokens for various services
}

// LoadConfig loads the configuration from the specified YAML file
func LoadConfig() (Config, error) {
	var cfg Config

	// Find the user's home directory
	currentUser, err := user.Current()
	if err != nil {
		return cfg, err
	}
	homeDir := currentUser.HomeDir

	viper.SetConfigType("yaml")
	// Set the name of the configuration file (without extension)
	viper.SetConfigName(".govcms-cli")
	// Set the path to look for the configuration file
	viper.AddConfigPath(homeDir)

	// Read in the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// If the config file is not found, create default config
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(homeDir, &cfg); err != nil {
				return cfg, err
			}
		} else {
			// If there's an error other than file not found, return it
			return cfg, err
		}
	}

	// Unmarshal the config file into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Creates and writes the default configuration to file
func createDefaultConfig(homeDir string, cfg *Config) error {
	// Set default values
	cfg.Tokens = map[string]string{
		"Github": "---",
		"Gitlab": "---",
		"Docker": "---",
	}

	// Set default workspace directory
	cfg.Workspace = filepath.Join(homeDir, "govcms")

	// Set the default tokens and workspace in the viper instance
	viper.Set("tokens", cfg.Tokens)
	viper.Set("workspace", cfg.Workspace)
	viper.Set("os", utils.GetOperatingSystem())

	// Write the configuration to file
	configFilePath := filepath.Join(homeDir, ".govcms-cli.yaml")
	if err := viper.WriteConfigAs(configFilePath); err != nil {
		return err
	}

	return nil
}
