package config

import "os"

type MySQLConfig struct {
	Address  string
	Username string
	Password string
	Database string
}

func LoadMySQLConfig() MySQLConfig {
	database := os.Getenv("MYSQL_DATABASE")
	if database == "" {
		database = "xiaoyaji"
	}

	return MySQLConfig{
		Address:  os.Getenv("MYSQL_ADDRESS"),
		Username: os.Getenv("MYSQL_USERNAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: database,
	}
}
