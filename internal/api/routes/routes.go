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
	JWT                jwt.Jwt
}

func SetupRoutes(cfg Config) {

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

	// Campaign Routes
	campaignGroup := cfg.Router.Group("/campaign")
	campaignGroup.Use(middlewares.Auth(cfg.JWT))
	{
		campaignGroup.POST("/create", cfg.CampaignHandler.HandleCreateCampaign)

		protected := campaignGroup.Use(middlewares.CampaignKey())
		{
			protected.GET("/:id", cfg.CampaignHandler.HandleGetCampaignByID)
			protected.PATCH("/:id", cfg.CampaignHandler.HandleUpdateCampaignByID)
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

	// Suggestions Routes
	suggestionsGroup := cfg.Router.Group("/suggestions")
	activitySuggestions := suggestionsGroup.Group("/activity")
	{
		activitySuggestions.GET("/:campaignID", cfg.SuggestionHandler.HandleGetActivitySuggestions)
		activitySuggestions.POST("/", cfg.SuggestionHandler.HandleGetActivitySuggestionsViaText)
	}

}
