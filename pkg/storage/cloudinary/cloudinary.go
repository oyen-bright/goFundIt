package cloudinary

import (
	"context"
	"strings"

	client "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/oyen-bright/goFundIt/pkg/storage"
)

type cloudinary struct {
	cloudinary *client.Cloudinary
	URL        string
}

func (c *cloudinary) UploadFile(filePath, folderPath string) (string, string, error) {
	ctx := context.Background()

	// Ensure the folder path is correctly formatted
	if folderPath != "" && !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	// Upload the file to Cloudinary
	uploadResult, err := c.cloudinary.Upload.Upload(ctx, filePath, uploader.UploadParams{
		Folder: folderPath,
	})
	if err != nil {
		return "", "", err
	}

	// Return the URL and the public ID of the uploaded file
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func (c *cloudinary) DeleteFile(publicID string) error {
	ctx := context.Background()

	// Delete the file from Cloudinary using its public ID
	_, err := c.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return err
	}

	return nil
}

func NewCloudinary(cloudinaryURL string) (storage.Storage, error) {

	cld, err := client.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}

	return &cloudinary{
		cloudinary: cld,
		URL:        cloudinaryURL,
	}, nil
}
