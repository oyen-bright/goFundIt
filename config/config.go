package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/oyen-bright/goFundIt/config/environment"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

type AppConfig struct {
	Environment      environment.Environment
	EmailProvider    providers.EmailProvider
	Port             string
	EmailConfig      email.EmailConfig
	EncryptionKey    []string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	XAPIKey          string
	JWTSecret        string
}

var BaseDir string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get working directory: %v", err))
	}
	BaseDir = filepath.Join(wd, "../../")
}

// LoadConfig loads the configuration for the application based on the environment.
//   - It parses the environment flag.
//   - Initializes the environment.
//   - Loads the corresponding .env file.
//   - If no flag for env is provided, it defaults to "" which is Development.
//   - To set env flag, run the application with -env=stg or -env=prod.
//   - Returns a Config struct with the environment and default port.
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

	envPath := filepath.Join(BaseDir, "config", "env", ".env."+environment.String())
	envData, err := loadEnv(envPath)
	if err != nil {
		return nil, err
	}

	emailPort, err := strconv.Atoi(envData["EMAIL_PORT"])
	if err != nil {
		return nil, err
	}

	postgresPort, err := strconv.Atoi(envData["POSTGRES_PORT"])
	if err != nil {
		return nil, err
	}

	emailProvider.Email(envData["EMAIL_PROVIDER"])

	emailConfig := email.EmailConfig{
		Host:           envData["EMAIL_HOST"],
		Port:           emailPort,
		From:           envData["EMAIL_NAME"],
		Username:       envData["EMAIL_USERNAME"],
		Password:       envData["EMAIL_PASSWORD"],
		SendGridAPIKey: envData["SENDGRID_API_KEY"],
	}

	return &AppConfig{
		Environment:      environment,
		Port:             envData["PORT"],
		EmailProvider:    emailProvider,
		EmailConfig:      emailConfig,
		EncryptionKey:    strings.Split(envData["ENCRYPTION_KEYS"], ","),
		PostgresDB:       envData["POSTGRES_DB"],
		PostgresUser:     envData["POSTGRES_USER"],
		PostgresPassword: envData["POSTGRES_PASSWORD"],
		PostgresHost:     envData["POSTGRES_HOST"],
		PostgresPort:     postgresPort,
		XAPIKey:          envData["X_API_KEY"],
		JWTSecret:        envData["JWT_SECRET"],
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
	requiredEnvs := []string{
		"PORT", "EMAIL_PROVIDER", "EMAIL_HOST", "EMAIL_PORT", "EMAIL_USERNAME",
		"EMAIL_PASSWORD", "ENCRYPTION_KEYS", "POSTGRES_DB", "POSTGRES_USER",
		"POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "X_API_KEY",
		"JWT_SECRET",
	}

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
