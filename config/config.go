package config

type Config struct {
	DB database
}

type database struct {
	Server     string
	Port       string
	Database   string
	Collection string
}
