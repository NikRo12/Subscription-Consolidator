package store

type Config struct {
	Driver      string `toml:"driver"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
