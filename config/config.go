package config

import (
	"flag"
	"fmt"

	"github.com/oyen-bright/goFundIt/config/environment"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Environment                    environment.Environment
	EmailProvider                  providers.EmailProvider
	FirebaseServiceAccountFilePath string `mapstructure:"firebase_service_account_file_path"`
	ServerPort                     string `mapstructure:"port"`
	GeminiKey                      string `mapstructure:"gemini_key"`
	PaystackKey                    string `mapstructure:"paystack_key"`
	EmailConfig                    email.EmailConfig
	CloudinaryURL                  string `mapstructure:"cloudinary_url"`
	AnalyticsReportEmail           string `mapstructure:"analytics_report_email"`
	DBConfig                       database.Config
	EncryptionKeys                 []string `mapstructure:"encryption_keys"` // Changed from string to []string
	XAPIKey                        string   `mapstructure:"x_api_key"`
	JWTSecret                      string   `mapstructure:"jwt_secret"`
}

type EmailConfigYAML struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	SendgridAPIKey string `mapstructure:"sendgrid_api_key"`
}

type DatabaseConfigYAML struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func LoadConfig() (*AppConfig, error) {
	env := parseEnvFlag()
	environment := environment.NewEnvironmentConfig(env).Current

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(fmt.Sprintf("config.%s", env))

	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var emailCfg EmailConfigYAML
	if err := v.UnmarshalKey("email", &emailCfg); err != nil {
		return nil, fmt.Errorf("failed to parse email config: %w", err)
	}

	var dbCfg DatabaseConfigYAML
	if err := v.UnmarshalKey("database", &dbCfg); err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.Environment = environment
	config.EmailProvider = providers.NewEmailProvider(providers.SMTP)
	config.EmailConfig = email.EmailConfig{
		Host:           emailCfg.Host,
		Port:           emailCfg.Port,
		From:           emailCfg.Name,
		Username:       emailCfg.Username,
		Password:       emailCfg.Password,
		SendGridAPIKey: emailCfg.SendgridAPIKey,
	}
	config.DBConfig = database.Config{
		Host:     dbCfg.Host,
		Port:     dbCfg.Port,
		Password: dbCfg.Password,
		DBName:   dbCfg.Name,
		User:     dbCfg.User,
	}

	return &config, nil
}

func parseEnvFlag() string {
	envFlag := flag.String("env", "dev", "Environment the application is running in")
	flag.Parse()
	return *envFlag
}
