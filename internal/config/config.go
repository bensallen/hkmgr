package config

import (
	"path/filepath"
)

// Config represents a hkmgr.toml config file
type Config struct {
	Network Network `toml:"network"`
	VM      VM      `toml:"vm"`
	Path    string  // Path to the loaded configuration
}

// UpdateRelativePaths finds relative paths in the config and turns them into
// fully qualified paths based on the config file path.
func (c *Config) UpdateRelativePaths() {
	configDir := filepath.Dir(c.Path)

	for name := range c.VM {
		c.VM[name].updateRelativePaths(configDir, name)
	}
}

// Defaults sets default values for unset variables in the config.
func (c *Config) Defaults() error {
	configDir := filepath.Dir(c.Path)

	for name := range c.VM {
		if err := c.VM[name].defaults(configDir, name); err != nil {
			return err
		}
	}
	return nil
}
