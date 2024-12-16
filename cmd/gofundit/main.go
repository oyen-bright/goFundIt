package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/internal/activity"
	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/campaign"
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"github.com/oyen-bright/goFundIt/internal/database"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/logger"
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

	tx := db.Debug()
	fmt.Println(tx.Statement)

	return cfg, db
}

func main() {

	//Initialize the logger
	logger := logger.New()

	// Initialize database and application configurations
	cfg, db := initialize()
	defer database.Close(db)

	// Create a new Gin router
	router := gin.Default()

	// Register general middlewares
	router.Use(middlewares.APIKey(cfg.XAPIKey))

	// Initialize encryption and email services
	encryptor := encryption.New(cfg.EncryptionKey)
	emailer := email.New(providers.EmailSMTP, cfg.EmailConfig)

	// Initialize jwt
	jwt := jwt.New(cfg.JWTSecret)

	// Initialize repositories
	authRepo := auth.Repository(db)
	otpRepo := otp.Repository(db)
	campaignRepo := campaign.Repository(db)
	contributionRepo := contributor.Repository(db)
	activityRepo := activity.Repository(db)

	// Create service instances
	otpService := otp.Service(otpRepo, emailer, *encryptor, logger)
	contributionService := contributor.Service(contributionRepo, logger)
	authService := auth.Service(authRepo, otpService, *encryptor, jwt, logger)
	campaignService := campaign.Service(campaignRepo, contributionService, authService, logger)
	activityService := activity.Service(activityRepo, authService, campaignService, logger)

	// Create handler instances
	authHandler := auth.Handler(authService)
	campaignHandler := campaign.Handler(campaignService)
	activityHandler := activity.Handler(activityService)

	// Register routes
	authHandler.RegisterRoutes(router.Group("/auth"), []gin.HandlerFunc{})
	campaignHandler.RegisterRoutes(router.Group("/campaign"), []gin.HandlerFunc{middlewares.Auth(jwt), middlewares.CampaignKey()})
	activityHandler.RegisterRoutes(router.Group("/activity"), []gin.HandlerFunc{middlewares.Auth(jwt), middlewares.CampaignKey()})
	// Start the server on the specified port
	router.Run(cfg.Port)
}
