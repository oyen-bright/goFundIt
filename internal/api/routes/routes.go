package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/internal/api/middlewares"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

type Config struct {
	Router          *gin.Engine
	AuthHandler     *handlers.AuthHandler
	CampaignHandler *handlers.CampaignHandler
	ActivityHandler *handlers.ActivityHandler
	JWT             jwt.Jwt
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
			protected.GET("/", cfg.CampaignHandler.HandleGetCampaigns)
			protected.GET("/:id", cfg.CampaignHandler.HandleGetCampaigns)
			protected.PATCH("/:id", cfg.CampaignHandler.HandleUpdateCampaign)
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
	}
}
