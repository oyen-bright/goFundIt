package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/auth"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// @Summary Authenticate User
// @Description Initiates the authentication process for a user by sending a verification code to their email address. Requires user's name and email in request body.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.User true "Auth Credentials"
// @Success 200 {object} SuccessResponse{data=models.Otp} "Please check your email for the OTP."
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs, please check and try again"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /auth [post]
func (a *AuthHandler) HandleAuth(context *gin.Context) {

	//TODO: use dto.AuthRequest
	var user models.User

	if err := context.BindJSON(&user); err != nil {
		BadRequest(context, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return
	}

	otp, err := a.service.RequestAuth(user.Email, *user.Name)
	if err != nil {
		FromError(context, err)
		return
	}

	Success(context, "Please check your email for the OTP.", otp.ToJSON())
}

// @Summary Verify OTP
// @Description Validates the verification code sent to user's email. Requires email, verification code, and request ID for verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyAuthRequest true "Verification Data"
// @Success 200 {object} SuccessResponse{data=dto.VerifyAuthResponse} "Authenticated"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs, please check and try again"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /auth/verify [post]
func (a *AuthHandler) HandleVerifyAuth(context *gin.Context) {
	//TODO: use dto.verifyAuthRequest
	var _otp models.Otp

	//Validate Request
	//Email- required, email
	//Code- required, string
	//RequestId- required, string
	if err := context.BindJSON(&_otp); err != nil {
		BadRequest(context, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return
	}

	//Verify Auth
	token, err := a.service.VerifyAuth(_otp.Email, _otp.Code, _otp.RequestId)
	if err != nil {
		FromError(context, err)
		return
	}

	Success(context, "Authenticated", dto.VerifyAuthResponse{
		Token: token,
	})

}

// @Summary Save FCM Token
// @Description Saves the FCM token for push notifications
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param request body dto.FCMUpdateRequest true "FCM Token"
// @Success 200 {object} SuccessResponse "FCM token saved"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs, please check and try again"
// @Router /fcm/save-token [post]
func (a *AuthHandler) HandleSaveFCMToken(c *gin.Context) {
	userHandle := getClaimsFromContext(c).Handle

	var req dto.FCMUpdateRequest

	//Validate Request
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return
	}

	//Save FCM Token
	err := a.service.SaveFCMToken(userHandle, req.FCMToken)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "FCM token saved", nil)
}
