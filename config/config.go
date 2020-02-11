package config

type Config struct {
	Database database1
}

type database1 struct {
	Server     string
	Port       string
	Database   string
	Collection string
}
