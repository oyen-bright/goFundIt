package handlers

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

// Helper functions

// getClaimsFromContext extracts JWT claims from the context
func getClaimsFromContext(c *gin.Context) jwt.Claims {
	return c.MustGet("claims").(jwt.Claims)
}

// getCampaignKey extracts the campaign key form the request
func getCampaignKey(c *gin.Context) string {
	return c.MustGet("Campaign-Key").(string)
}

// GetCampaignID extracts the campaignId form the requestParam
func GetCampaignID(c *gin.Context) string {
	return c.Param("campaignID")
}

// GetCommentID extracts the commentID form the requestParam
func getCommentID(c *gin.Context) string {
	return c.Param("commentID")
}

// parseActivityID converts the activity ID from the URL parameter to uint
func parseActivityID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("activityID"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// parseContributorID converts the contributor ID from the URL parameter to uint
func parseContributorID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("contributorID"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateTempFileFromMultipart creates a temporary file from multipart form data
func createTempFileFromMultipart(file *multipart.FileHeader) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "upload-*.png")
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if _, err = io.Copy(tempFile, src); err != nil {
		return nil, err
	}

	if _, err = tempFile.Seek(0, 0); err != nil {
		return nil, err
	}

	return tempFile, nil
}
