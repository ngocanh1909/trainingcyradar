package config

type Config struct {
	DB struct {
		Server     string `toml:"server"`
		Database   string `toml:"database"`
		Collection string `toml:"collection"`
		Username   string `toml:"username"`
		Password   string `toml:"password"`
	} `toml:"db"`
}
