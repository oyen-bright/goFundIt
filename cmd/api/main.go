package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/providers"

	"github.com/oyen-bright/goFundIt/internal/ai/gemini"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/internal/api/middlewares"
	"github.com/oyen-bright/goFundIt/internal/api/routes"
	postgress "github.com/oyen-bright/goFundIt/internal/repositories/postgres"
	"github.com/oyen-bright/goFundIt/internal/services"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/oyen-bright/goFundIt/pkg/websocket"

	"gorm.io/gorm"
)

func initialize() (*config.AppConfig, *gorm.DB) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.Init(cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	db.Debug()

	return cfg, db
}

func main() {

	//Initialize the logger
	logger := logger.New()

	// Initialize database and application configurations
	cfg, db := initialize()
	defer database.Close(db)

	// Initialize Core Services
	aiClient, _ := gemini.NewClient(cfg.GeminiKey)
	defer gemini.Close(aiClient)

	encryptor := encryption.New(cfg.EncryptionKey)
	emailer := email.New(providers.EmailSMTP, cfg.EmailConfig)
	jwtService := jwt.New(cfg.JWTSecret)

	websocketHub := websocket.NewHub()
	go websocketHub.Run()
	defer websocketHub.Close()

	// Initialize Repositories
	authRepo := postgress.NewAuthRepository(db)
	otpRepo := postgress.NewOTPRepository(db)
	campaignRepo := postgress.NewCampaignRepository(db)
	contributorRepo := postgress.NewContributorRepository(db)
	activityRepo := postgress.NewActivityRepo(db)
	commentRepo := postgress.NewCommentRepository(db)

	// initialize the event broadcaster
	eventBroadcaster := services.NewEventBroadcaster(websocketHub)

	// Initialize Services
	otpService := services.NewOTPService(otpRepo, emailer, *encryptor, logger)
	authService := services.NewAuthService(authRepo, otpService, *encryptor, jwtService, logger)
	campaignService := services.NewCampaignService(campaignRepo, authService, eventBroadcaster, logger)
	contributorService := services.NewContributorService(contributorRepo, campaignService, eventBroadcaster, logger)
	activityService := services.NewActivityService(activityRepo, authService, campaignService, eventBroadcaster, logger)
	commentService := services.NewCommentService(commentRepo, authService, activityService, eventBroadcaster, logger)
	suggestionService := services.NewSuggestionService(aiClient, campaignService, logger)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	activityHandler := handlers.NewActivityHandler(activityService)
	contributorHandler := handlers.NewContributorHandler(contributorService)
	commentHandler := handlers.NewCommentHandler(commentService)
	suggestionHandler := handlers.NewSuggestionHandler(suggestionService)
	websocketHandler := handlers.NewWebSocketHandler(websocketHub, campaignService)

	// Initialize Gin Router
	router := gin.Default()
	router.Use(middlewares.APIKey(cfg.XAPIKey))

	// Setup Routes
	routes.SetupRoutes(routes.Config{
		Router:             router,
		AuthHandler:        authHandler,
		CampaignHandler:    campaignHandler,
		ContributorHandler: contributorHandler,
		ActivityHandler:    activityHandler,
		CommentHandler:     commentHandler,
		SuggestionHandler:  suggestionHandler,
		WebSocketHandler:   websocketHandler,
		JWT:                jwtService,
	})

	// Start Server
	if err := router.Run(cfg.ServerPort); err != nil {
		panic("Failed to start server:")
	}
}
