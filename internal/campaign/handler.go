package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/interfaces"
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	"github.com/oyen-bright/goFundIt/internal/utils/response"
)

type campaignHandler struct {
	service CampaignService
}

func Handler(service CampaignService) interfaces.HandlerInterface {
	return &campaignHandler{
		service: service,
	}
}

// RegisterRoutes registers the campaign-related routes with the provided RouterGroup.
// It applies the provided middlewares to the routes.
//
// Routes:
//   - POST /create: Creates a new campaign.
//   - GET /: Retrieves all campaigns.
//   - GET /:id: Retrieves a specific campaign by ID.
//   - PATCH /:id: Updates a specific campaign by ID.
func (c *campaignHandler) RegisterRoutes(campaignRoute *gin.RouterGroup, middlewares []gin.HandlerFunc) {
	campaignRoute.Use(middlewares[0])
	campaignRoute.POST("/create", c.handleCreateCampaign)

	protectedCampaignRoute := campaignRoute.Group("/")
	protectedCampaignRoute.Use(middlewares[1])
	protectedCampaignRoute.GET("/", c.handleGetCampaigns)
	protectedCampaignRoute.GET("/:id", c.handleGetCampaigns)
	protectedCampaignRoute.PATCH("/:id", c.handleUpdateCampaign)

}

// handleCreateCampaign handles incoming campaign creation requests.
// Required Fields:
//   - Title: required, string
//   - Description: optional, string
//   - Images: []{imageUrl: required, url}
//   - Activities: []{title: required, string, subtitle: optional, string, imageUrl: optional, string, isMandatory: optional, bool}
//   - TargetAmount: required, number
//   - StartDate: required, date
//   - EndDate: required, date
func (c *campaignHandler) handleCreateCampaign(context *gin.Context) {

	claims := context.MustGet("claims").(jwt.Claims)
	var campaign Campaign

	//bind request to the campaign model
	if err := context.BindJSON(&campaign); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", response.ExtractValidationErrors(err))
		return
	}
	//validate the request
	if err := campaign.ValidateNewCampaign(); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", response.ExtractValidationErrors(err))
		return
	}

	//create service
	campaign, err := c.service.CreateCampaign(campaign, claims.Handle)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Campaigns created successfully", campaign.ToJSON())

}

func (c *campaignHandler) handleGetCampaigns(context *gin.Context) {

	response.Success(context, "Campaigns retrieved successfully", []string{})

}

func (c *campaignHandler) handleUpdateCampaign(context *gin.Context) {

	response.Success(context, "Campaigns Updated successfully", []string{})

}
