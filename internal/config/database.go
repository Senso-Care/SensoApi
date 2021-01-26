package config

type DatabaseConfiguration struct {
	ConnectionUri   string
	DbName          string
	RetentionPolicy string
	Username        string
	Password        string
}
