package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/internal/api/middlewares"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

type Config struct {
	Router             *gin.Engine
	AuthHandler        *handlers.AuthHandler
	CampaignHandler    *handlers.CampaignHandler
	SuggestionHandler  *handlers.SuggestionHandler
	ContributorHandler *handlers.ContributorHandler
	CommentHandler     *handlers.CommentHandler
	ActivityHandler    *handlers.ActivityHandler
	WebSocketHandler   *handlers.WebSocketHandler
	PaymentHandler     *handlers.PaymentHandler
	PayoutHandler      *handlers.PayoutHandler
	AnalyticsHandler   *handlers.AnalyticsHandle
	PaystackKey        string
	XAPIKey            string
	JWT                jwt.Jwt
}

func SetupRoutes(cfg Config) {

	// Add Sentry middleware
	// cfg.Router.Use(sentrygin.New(sentrygin.Options{
	// 	Repanic:         true,
	// 	WaitForDelivery: true,
	// 	Timeout:         2 * time.Second,
	// }))

	// Webhook route
	cfg.Router.POST("/payment/paystack/webhook", cfg.PaymentHandler.HandlePayStackWebhook, middlewares.PaystackSignature(cfg.PaystackKey))

	// API Key Middleware
	cfg.Router.Use(middlewares.APIKey(cfg.XAPIKey))

	// Websocket Routes
	ws := cfg.Router.Group("/ws")
	ws.Use(middlewares.Auth(cfg.JWT), middlewares.CampaignKey())
	{
		ws.GET("/campaign/:campaignID", cfg.WebSocketHandler.HandleCampaignWebSocket)
	}

	// Auth Routes
	authGroup := cfg.Router.Group("/auth")
	{
		authGroup.POST("/", cfg.AuthHandler.HandleAuth)
		authGroup.POST("/verify", cfg.AuthHandler.HandleVerifyAuth)
	}

	// FCM Routes
	fcmGroup := cfg.Router.Group("/fcm")
	fcmGroup.Use(middlewares.Auth(cfg.JWT))
	{
		fcmGroup.POST("/save-token", cfg.AuthHandler.HandleSaveFCMToken)
		//TODO: Implement the route handler
		// fcmGroup.DELETE("/remove-token", cfg.AuthHandler.HandleRemoveToken)
	}

	// Campaign Routes
	campaignGroup := cfg.Router.Group("/campaign")
	campaignGroup.Use(middlewares.Auth(cfg.JWT))
	{
		campaignGroup.POST("/create", cfg.CampaignHandler.HandleCreateCampaign)

		protected := campaignGroup.Use(middlewares.CampaignKey())
		{
			protected.GET("/:campaignID", cfg.CampaignHandler.HandleGetCampaignByID)
			protected.PATCH("/:campaignID", cfg.CampaignHandler.HandleUpdateCampaignByID)
		}
	}

	// Activity Routes
	activityGroup := cfg.Router.Group("/activity")
	activityGroup.Use(middlewares.Auth(cfg.JWT), middlewares.CampaignKey())
	{
		activityGroup.GET("/:campaignID", cfg.ActivityHandler.HandleGetActivitiesByCampaignID)
		activityGroup.GET("/:campaignID/:activityID", cfg.ActivityHandler.HandleGetActivityByID)
		activityGroup.POST("/:campaignID", cfg.ActivityHandler.HandleCreateActivity)
		activityGroup.PATCH("/:campaignID/:activityID", cfg.ActivityHandler.HandleUpdateActivity)
		activityGroup.DELETE("/:campaignID/:activityID", cfg.ActivityHandler.HandleDeleteActivityByID)

		activityGroup.POST("/:campaignID/:activityID/approve", cfg.ActivityHandler.HandleApproveActivity)

		participation := activityGroup.Group("/:campaignID/:activityID/participants")
		{
			participation.POST("/:contributorID", cfg.ActivityHandler.HandleOptInContributor)
			participation.DELETE("/:contributorID", cfg.ActivityHandler.HandleOptOutContributor)
			participation.GET("/", cfg.ActivityHandler.HandleGetParticipants)
		}

		comments := activityGroup.Group("/:campaignID/:activityID/comments")
		{
			comments.POST("/", cfg.CommentHandler.HandleCreateComment)
			comments.PATCH("/:commentID", cfg.CommentHandler.HandleUpdateComment)
			comments.GET("/", cfg.CommentHandler.HandleGetActivityComments)
			comments.GET("/:commentID/replies", cfg.CommentHandler.HandleGetCommentReplies)
			comments.DELETE("/:commentID", cfg.CommentHandler.HandleDeleteComment)
		}
	}

	// Contributor Routes
	contributorGroup := cfg.Router.Group("/contributor")
	contributorGroup.Use(middlewares.Auth(cfg.JWT), middlewares.CampaignKey())
	{
		contributorGroup.POST("/:campaignID", cfg.ContributorHandler.HandleAddContributor)
		contributorGroup.DELETE("/:campaignID/:contributorID", cfg.ContributorHandler.HandleRemoveContributor)
		contributorGroup.PATCH("/:campaignID/:contributorID", cfg.ContributorHandler.HandleEditContributor)
		contributorGroup.GET("/:campaignID", cfg.ContributorHandler.HandleGetContributorsByCampaignID)
		contributorGroup.GET("/:campaignID/:contributorID", cfg.ContributorHandler.HandleGetContributorByID)
	}

	// Payment Routes
	paymentGroup := cfg.Router.Group("/payment")
	paymentGroup.Use(middlewares.Auth(cfg.JWT), middlewares.CampaignKey())
	{
		//TODO: fix route name
		paymentGroup.POST("/contributor/:contributorID", cfg.PaymentHandler.HandleInitializePayment)
		paymentGroup.POST("manual/contributor/:contributorID", cfg.PaymentHandler.HandleInitializeManualPayment)
		// Payment verification route
		paymentGroup.POST("/verify/:reference", cfg.PaymentHandler.HandleVerifyPayment)
		paymentGroup.POST("/manual/verify/:reference", cfg.PaymentHandler.HandleVerifyManualPayment)
	}

	// Payout Routes
	payoutGroup := cfg.Router.Group("/payout")
	{
		payoutGroup.GET("/bank-list", cfg.PayoutHandler.HandleGetBankList)
		payoutGroup.POST("/verify/bank-account", cfg.PayoutHandler.HandleVerifyAccount)
	}
	payoutGroup.Use(middlewares.Auth(cfg.JWT), middlewares.CampaignKey())
	{
		payoutGroup.POST("/:campaignID", cfg.PayoutHandler.HandleInitializePayout)
		payoutGroup.POST("manual/:campaignID", cfg.PayoutHandler.HandleInitializeManualPayout)
		payoutGroup.GET("/:campaignID", cfg.PayoutHandler.HandleGetPayoutByCampaignID)

	}

	// Suggestions Routes
	suggestionsGroup := cfg.Router.Group("/suggestions")
	activitySuggestions := suggestionsGroup.Group("/activity")
	{
		activitySuggestions.POST("/", cfg.SuggestionHandler.HandleGetActivitySuggestionsViaText)

		activitySuggestions.Use(middlewares.CampaignKey())
		activitySuggestions.GET("/:campaignID", cfg.SuggestionHandler.HandleGetActivitySuggestions)
	}

	//Analytics routes
	analyticsGroup := cfg.Router.Group("/analytics")
	{
		analyticsGroup.GET("/process", cfg.AnalyticsHandler.HandleProcessAnalyticsNow)
	}

}
