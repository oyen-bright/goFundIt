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
	ContributorHandler *handlers.ContributorHandler
	ActivityHandler    *handlers.ActivityHandler
	JWT                jwt.Jwt
}

func SetupRoutes(cfg Config) {
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
			participation.POST("/:contributorID", cfg.ActivityHandler.HandleOptInContributor)    // Opt in
			participation.DELETE("/:contributorID", cfg.ActivityHandler.HandleOptOutContributor) // Opt out
			participation.GET("/", cfg.ActivityHandler.HandleGetParticipants)                    // List participants
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

}
