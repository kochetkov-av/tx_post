package main

import "fmt"

type Config struct {
	DbConfig struct {
		Host     string
		Port     string
		Dbname   string
		User     string
		Password string
		Sslmode  string
	} `yaml:"db"`
	Sources              []string
	CancellationInterval int `yaml:"cancellation_interval"`
}

func (config *Config) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.DbConfig.Host,
		config.DbConfig.Port,
		config.DbConfig.User,
		config.DbConfig.Dbname,
		config.DbConfig.Password,
		config.DbConfig.Sslmode,
	)
}
