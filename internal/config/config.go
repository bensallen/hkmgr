package config

type Config struct {
	Network Network `toml:"network"`
	VM      VM      `toml:"vm"`
	Path    string  // Path to the loaded configuration
}

// UpdateRelativePaths finds relative paths in the config and turns them into
// fully qualified paths based on the config file path.
func (c *Config) UpdateRelativePaths() {
	for name := range c.VM {
		c.VM[name].updateRelativePaths(c.Path, name)
	}
}
