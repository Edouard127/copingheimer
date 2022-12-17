package intf

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	InstanceCount int    `json:"instance_count" env:"INSTANCE_COUNT"`
	Timeout       int    `json:"timeout" env:"TIMEOUT"`
	CPUSaver      bool   `json:"cpu_saver" env:"CPU_SAVER"`
	Database      string `json:"database" env:"DATABASE"` // The type of the database, currently supported: bolt, mongodb
	DatabaseURL   string `json:"database_url" env:"DATABASE_URL"`
}

func LoadConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var config Config
	if strings.Contains(filename, ".json") {
		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			panic(err)
		}
	} else if strings.Contains(filename, ".env") {
		// Load the environment variables
		if err := godotenv.Load(file.Name()); err != nil {
			panic(err)
		}
		err = env.Parse(&config)
		if err != nil {
			panic(err)
		}
	}

	return &config
}
