package config

import (
	"errors"
	"flag"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/oyen-bright/goFundIt/config/environment"
	providers "github.com/oyen-bright/goFundIt/config/provider"
)

type AppConfig struct {
	Environment    environment.Environment
	EmailProvider  providers.EmailProvider
	Port           int
	EmailHost      string
	EmailPort      int
	EmailUsername  string
	EmailPassword  string
	EmailName      string
	SendGridAPIKey string
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
func LoadConfig() (*AppConfig, error) {

	var env string
	var environment environment.Environment
	var emailProvider providers.EmailProvider

	envFlag := flag.String("env", "", "Environment the application is running in")
	flag.Parse()

	env = *envFlag
	environment.New(env)

	envPath := ".env." + environment.String()
	envData, err := loadEnv(envPath)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(envData["PORT"])
	if err != nil {
		return nil, err
	}

	emailPort, err := strconv.Atoi(envData["EMAIL_PORT"])
	if err != nil {
		return nil, err
	}
	emailProvider.Email(envData["EMAIL_PROVIDER"])

	return &AppConfig{
		Environment:    environment,
		Port:           port,
		EmailProvider:  emailProvider,
		EmailHost:      envData["EMAIL_HOST"],
		EmailPort:      emailPort,
		EmailUsername:  envData["EMAIL_USERNAME"],
		EmailPassword:  envData["EMAIL_PASSWORD"],
		EmailName:      envData["EMAIL_NAME"],
		SendGridAPIKey: envData["SENDGRID_API_KEY"],
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
	requiredEnvs := []string{"PORT", "EMAIL_PROVIDER", "EMAIL_HOST", "EMAIL_PORT", "EMAIL_USERNAME", "EMAIL_PASSWORD"}

	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	envData, err := godotenv.Read(envPath)
	if err != nil {
		return nil, err
	}
	for _, env := range requiredEnvs {
		if _, isAvailable := envData[env]; !isAvailable {
			return nil, errors.New("Missing required environment variable: " + env)
		}
	}
	return envData, nil
}
