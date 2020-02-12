package config

type Config struct {
	DB database
}

type database struct {
	Server     string
	Database   string
	Collection string
	Username   string
	Password   string
}
