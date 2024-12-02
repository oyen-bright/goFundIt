package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment Environment
	Port        int
}

// LoadConfig loads the configuration for the application based on the environment.
// - It parses the environment flag.
// - Initializes the environment.
// - Loads the corresponding .env file.
// - If no flag for env is provided, it defaults to "" which is Development.
// - To set env flag, run the application with -env=stg or -env=prod.
// - Returns a Config struct with the environment and default port.
//
// Returns:
//   - *Config: A pointer to the Config struct containing the environment and port.
//   - error: An error if any occurs during the loading of the configuration.
func LoadConfig() (*Config, error) {

	var env string
	var environment Environment

	envFlag := flag.String("env", "", "Environment the application is running in")
	flag.Parse()

	env = *envFlag
	environment.init(env)

	envPath := ".env." + environment.String()
	envData, err := loadEnv(envPath)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(envData["PORT"])
	if err != nil {
		return nil, err
	}

	return &Config{
		Environment: environment,
		Port:        port,
	}, nil
}

// loadEnv loads environment variables from a specified file and ensures that
// all required environment variables are present.
//
// Parameters:
//   - envPath: The path to the environment file to load.
//
// Returns:
//   - A map containing the environment variables and their values.
//   - An error if there is an issue loading the environment file or if any
//     required environment variables are missing.
//
// Required Environment Variables:
//   - PORT: The port number on which the application should run.
func loadEnv(envPath string) (map[string]string, error) {
	requiredEnvs := []string{"PORT"}

	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	envData, err := godotenv.Read(envPath)
	if err != nil {
		return nil, err
	}
	for _, env := range requiredEnvs {
		fmt.Println(env)
		if _, isAvailable := envData[env]; !isAvailable {
			return nil, errors.New("Missing required environment variable: " + env)
		}
	}
	return envData, nil
}
