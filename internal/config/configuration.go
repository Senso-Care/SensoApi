package config

type Configuration struct {
	Database DatabaseConfiguration
	Mock     bool
	Cores    int
	Port     int
}
