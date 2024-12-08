package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/database"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/middlewares"
	"gorm.io/gorm"
)

func initialize() (*config.AppConfig, *gorm.DB) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.Init(*cfg)
	if err != nil {
		panic(err)
	}
	return cfg, db
}

func main() {

	// Initialize database and application configurations
	cfg, db := initialize()
	// migrations.DropOtpTable(db)
	defer database.Close(db)

	// Create a new Gin router
	r := gin.Default()

	// Register general middlewares
	r.Use(middlewares.APIKeyAuthMiddleware(cfg.XAPIKey))

	// Initialize encryption and email services
	encryptor := encryption.New(cfg.EncryptionKey)
	emailer := email.New(providers.EmailSMTP, cfg.EmailConfig)

	// Initialize jwt
	jwt := jwt.New(cfg.JWTSecret)

	// Initialize repositories
	authRepo := auth.Repository(db)
	otpRepo := otp.Repository(db)

	// Create service instances
	otpService := otp.Service(otpRepo, emailer, *encryptor)
	authService := auth.Service(authRepo, *encryptor, jwt)

	// Create handler instances
	authHandler := auth.Handler(otpService, authService)

	// Register routes
	authHandler.RegisterRoutes(r.Group("/auth"), []gin.HandlerFunc{})

	// Start the server on the specified port
	r.Run(cfg.Port)
}
