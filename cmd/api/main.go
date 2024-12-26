package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/internal/ai/gemini"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/internal/api/routes"
	postgress "github.com/oyen-bright/goFundIt/internal/repositories/postgres"
	"github.com/oyen-bright/goFundIt/internal/services"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/fcm"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	"github.com/oyen-bright/goFundIt/pkg/storage/cloudinary"
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
	// Initialize database and application configurations
	cfg, db := initialize()
	defer database.Close(db)

	if err := config.InitSentry(
		cfg.Environment.String(),
		cfg.Environment.IsDevelopment(),
	); err != nil {
		panic(err)
	}
	defer sentry.Flush(2 * time.Second)

	//Initialize the logger
	logger := logger.New()

	// Initialize AI Client
	aiClient, _ := gemini.NewClient(cfg.GeminiKey)
	defer gemini.Close(aiClient)

	// Initialize Websocket Hub
	websocketHub := websocket.NewHub()
	go websocketHub.Run()
	defer websocketHub.Close()

	storage, err := cloudinary.NewCloudinary(cfg.CloudinaryURL)
	if err != nil {
		panic(err)
	}

	//TODO: Initialize FCM Client
	fcmClient, err := fcm.New("")

	// Initialize Core Services

	encryptor := encryption.New(cfg.EncryptionKey)
	emailer := email.New(providers.EmailSMTP, cfg.EmailConfig)
	jwtService := jwt.New(cfg.JWTSecret)

	//Initialize paystack
	paystackClient := paystack.NewClient(cfg.PaystackKey)

	// Initialize Repositories
	authRepo := postgress.NewAuthRepository(db)
	otpRepo := postgress.NewOTPRepository(db)
	campaignRepo := postgress.NewCampaignRepository(db)
	contributorRepo := postgress.NewContributorRepository(db)
	activityRepo := postgress.NewActivityRepo(db)
	commentRepo := postgress.NewCommentRepository(db)
	paymentRepo := postgress.NewPaymentRepository(db)
	payoutRepo := postgress.NewPayoutRepository(db)
	analyticsRepo := postgress.NewAnalyticsRepository(db)

	// initialize the event broadcaster
	eventBroadcaster := services.NewEventBroadcaster(websocketHub)

	// Initialize Services
	otpService := services.NewOTPService(otpRepo, emailer, *encryptor, logger)
	authService := services.NewAuthService(authRepo, otpService, *encryptor, jwtService, logger)
	notificationService := services.NewNotificationService(emailer, authService, fcmClient, logger)
	campaignService := services.NewCampaignService(campaignRepo, authService, notificationService, eventBroadcaster, logger)
	contributorService := services.NewContributorService(contributorRepo, campaignService, notificationService, eventBroadcaster, logger)
	activityService := services.NewActivityService(activityRepo, authService, campaignService, eventBroadcaster, notificationService, logger)
	commentService := services.NewCommentService(commentRepo, authService, activityService, notificationService, eventBroadcaster, logger)
	suggestionService := services.NewSuggestionService(aiClient, campaignService, logger)
	paymentService := services.NewPaymentService(paymentRepo, contributorService, campaignService, notificationService, paystackClient, storage, eventBroadcaster, logger)
	payoutService := services.NewPayoutService(payoutRepo, campaignService, notificationService, paystackClient, eventBroadcaster, logger)

	analyticsService := services.NewAnalyticsService(campaignService, authService, analyticsRepo, cfg.AnalyticsReportEmail, emailer, logger)
	if err := analyticsService.StartAnalytics(); err != nil {
		panic(err)
	}
	defer analyticsService.StopAnalytics()

	cronService := services.NewCronService(campaignService, notificationService, logger)
	if err := cronService.StartCronJobs(); err != nil {
		panic(err)
	}
	defer cronService.StopCronJobs()

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	activityHandler := handlers.NewActivityHandler(activityService)
	contributorHandler := handlers.NewContributorHandler(contributorService)
	commentHandler := handlers.NewCommentHandler(commentService)
	suggestionHandler := handlers.NewSuggestionHandler(suggestionService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	payoutHandler := handlers.NewPayoutHandler(payoutService)
	websocketHandler := handlers.NewWebSocketHandler(websocketHub, campaignService)

	// Initialize Gin Router
	router := gin.Default()

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
		PaymentHandler:     paymentHandler,
		PayoutHandler:      payoutHandler,
		PaystackKey:        cfg.PaystackKey,
		XAPIKey:            cfg.XAPIKey,
		JWT:                jwtService,
	})

	// Start Server
	if err := router.Run(cfg.ServerPort); err != nil {
		panic("Failed to start server:")
	}
}
