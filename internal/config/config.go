package config

type Config struct {
	Network Network `toml:"network"`
	VM      VM      `toml:"vm"`
}
