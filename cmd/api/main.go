package main

// @title GoFundIt API
// @version 1.0
// @description GoFundIt is a collaborative platform enabling families and friends to plan and fund group activities, such as vacations and events. Features end-to-end encryption, secure payment handling, and flexible contribution options.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://gofundit.com/support
// @contact.email oyeniyibright@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key for accessing the API endpoints

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token for authenticated requests. Use the format: Bearer <token>

// @securityDefinitions.apikey CampaignKeyAuth
// @in header
// @name Campaign-Key
// @description Campaign Key required for campaign-specific operations

// @tag.name auth
// @tag.description Secure authentication with end-to-end encryption, including login and OTP verification

// @tag.name campaign
// @tag.description Group activity and event campaign management with encrypted data storage

// @tag.name contributor
// @tag.description Manage group members and their contributions with secure payment handling

// @tag.name activity
// @tag.description Track and manage group activities and events securely

// @tag.name payment
// @tag.description Secure payment processing with flexible contribution options

// @tag.name payout
// @tag.description Safe fund distribution for group activities

// @tag.name suggestion
// @tag.description AI-powered suggestions for group activities and event planning

// @tag.name comment
// @tag.description Encrypted group communication with comments and replies

// @tag.name analytics
// @tag.description Process analytics for platform performance

// @tag.name websocket
// @tag.description Real-time encrypted updates for group coordination

// @Security ApiKeyAuth

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	doc "github.com/oyen-bright/goFundIt/cmd/docs/swagger"
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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initialize() (*config.AppConfig, *gorm.DB) {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.Init(cfg.DBConfig, cfg.Environment.IsDevelopment())
	if err != nil {
		panic(err)
	}
	if cfg.Environment.IsDevelopment() {
		db.Debug()
	}
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
	logger := logger.New(cfg.Environment.IsProduction(), cfg.Environment.IsDevelopment())

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

	fcmClient, err := fcm.New(cfg.FirebaseServiceAccountFilePath)
	if err != nil {
		panic(err)
	}

	// Initialize Core Services

	encryptor := encryption.New(cfg.EncryptionKeys)
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

	analyticsService := services.NewAnalyticsService(analyticsRepo, cfg.AnalyticsReportEmail, emailer, logger)
	if err := analyticsService.StartAnalytics(); err != nil {
		panic(err)
	}
	defer analyticsService.StopAnalytics()

	otpService := services.NewOTPService(otpRepo, emailer, logger)
	authService := services.NewAuthService(authRepo, otpService, encryptor, analyticsService, jwtService, logger)
	notificationService := services.NewNotificationService(emailer, authService, fcmClient, logger)
	campaignService := services.NewCampaignService(campaignRepo, authService, analyticsService, notificationService, encryptor, eventBroadcaster, logger)
	contributorService := services.NewContributorService(contributorRepo, campaignService, analyticsService, authService, notificationService, eventBroadcaster, logger)
	activityService := services.NewActivityService(activityRepo, authService, campaignService, eventBroadcaster, analyticsService, notificationService, logger)
	commentService := services.NewCommentService(commentRepo, authService, activityService, notificationService, eventBroadcaster, logger)
	suggestionService := services.NewSuggestionService(aiClient, campaignService, logger)
	paymentService := services.NewPaymentService(paymentRepo, contributorService, analyticsService, campaignService, notificationService, paystackClient, storage, eventBroadcaster, logger)
	payoutService := services.NewPayoutService(payoutRepo, campaignService, notificationService, paystackClient, eventBroadcaster, logger)

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
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	websocketHandler := handlers.NewWebSocketHandler(websocketHub, campaignService)

	if cfg.Environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	// Initialize Gin Router
	router := gin.Default()

	// Configure Swagger
	doc.SwaggerInfo.BasePath = "/"

	// Redirect root to Swagger docs
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(1),
		ginSwagger.DocExpansion("list"),
		ginSwagger.PersistAuthorization(true),
	))

	// Setup Routes
	routes.SetupRoutes(routes.Config{
		Router:             router,
		AuthHandler:        authHandler,
		CampaignHandler:    campaignHandler,
		ContributorHandler: contributorHandler,
		ActivityHandler:    activityHandler,
		AnalyticsHandler:   analyticsHandler,
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
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}
	if err := router.Run(":" + port); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
