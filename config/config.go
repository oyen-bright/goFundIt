package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/oyen-bright/goFundIt/config/environment"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

type AppConfig struct {
	Environment          environment.Environment
	EmailProvider        providers.EmailProvider
	ServerPort           string
	GeminiKey            string
	PaystackKey          string
	EmailConfig          email.EmailConfig
	CloudinaryURL        string
	AnalyticsReportEmail string
	DBConfig             database.Config
	EncryptionKey        []string
	XAPIKey              string
	JWTSecret            string
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
func LoadConfig() (*AppConfig, error) {
	env := parseEnvFlag()
	environment := initializeEnvironment(env)

	envData, err := loadEnvFile(environment)
	if err != nil {
		return nil, err
	}

	appConfig, err := parseEnvData(envData)
	if err != nil {
		return nil, err
	}

	appConfig.Environment = environment

	return appConfig, nil
}

func parseEnvFlag() string {
	//Default env is dev
	envFlag := flag.String("env", "dev", "Environment the application is running in")
	flag.Parse()
	return *envFlag
}

func initializeEnvironment(env string) environment.Environment {
	return environment.NewEnvironmentConfig(env).Current
}

func loadEnvFile(env environment.Environment) (map[string]string, error) {
	envPath := filepath.Join(BaseDir, "config", "env", ".env."+env.String())
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}
	return godotenv.Read(envPath)
}

func parseEnvData(envData map[string]string) (*AppConfig, error) {
	requiredEnvs := []string{
		"PORT", "EMAIL_PROVIDER", "EMAIL_HOST", "EMAIL_PORT", "EMAIL_USERNAME",
		"EMAIL_PASSWORD", "ENCRYPTION_KEYS", "POSTGRES_DB", "POSTGRES_USER",
		"POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "X_API_KEY",
		"JWT_SECRET", "GEMINI_KEY", "PAYSTACK_KEY", "CLOUDINARY_URL", "ANALYTICS_REPORT_EMAIL",
	}

	if err := checkRequiredEnvs(envData, requiredEnvs); err != nil {
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
	var emailProvider providers.EmailProvider
	providers.NewEmailProvider(providers.SMTP)

	emailConfig := email.EmailConfig{
		Host:           envData["EMAIL_HOST"],
		Port:           emailPort,
		From:           envData["EMAIL_NAME"],
		Username:       envData["EMAIL_USERNAME"],
		Password:       envData["EMAIL_PASSWORD"],
		SendGridAPIKey: envData["SENDGRID_API_KEY"],
	}

	dbConfig := database.Config{
		Host:     envData["POSTGRES_HOST"],
		Port:     postgresPort,
		Password: envData["POSTGRES_PASSWORD"],
		DBName:   envData["POSTGRES_DB"],
		User:     envData["POSTGRES_USER"],
	}

	return &AppConfig{
		ServerPort:           envData["PORT"],
		GeminiKey:            envData["GEMINI_KEY"],
		PaystackKey:          envData["PAYSTACK_KEY"],
		AnalyticsReportEmail: envData["ANALYTICS_REPORT_EMAIL"],
		EmailProvider:        emailProvider,
		EmailConfig:          emailConfig,
		EncryptionKey:        strings.Split(envData["ENCRYPTION_KEYS"], ","),
		DBConfig:             dbConfig,
		XAPIKey:              envData["X_API_KEY"],
		CloudinaryURL:        envData["CLOUDINARY_URL"],
		JWTSecret:            envData["JWT_SECRET"],
	}, nil
}

func checkRequiredEnvs(envData map[string]string, requiredEnvs []string) error {
	for _, env := range requiredEnvs {
		if _, isAvailable := envData[env]; !isAvailable {
			return errors.New("Missing required environment variable: " + env)
		}
	}
	return nil
}
