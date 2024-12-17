package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a *AuthHandler) HandleAuth(context *gin.Context) {
	var user models.User

	//Validate Request
	//Email- required, email
	//Name- optional, string
	if err := context.BindJSON(&user); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	otp, err := a.service.RequestAuth(user.Email, user.Name)
	if err != nil {
		response.FromError(context, err)
		return
	}

	response.Success(context, "Please check your email for the OTP.", otp.ToJSON())

}

func (a *AuthHandler) HandleVerifyAuth(context *gin.Context) {
	var _otp models.Otp

	//Validate Request
	//Email- required, email
	//Code- required, string
	//RequestId- required, string
	if err := context.BindJSON(&_otp); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	//Verify Auth
	token, err := a.service.VerifyAuth(_otp.Email, _otp.Code, _otp.RequestId)
	if err != nil {
		response.FromError(context, err)
		return
	}

	response.Success(context, "Authenticated", gin.H{
		"token": token,
	})

}
