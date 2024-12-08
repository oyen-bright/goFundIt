package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/interfaces"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/utils/response"
)

type authHandler struct {
	otpService  otp.OTPService
	authService AuthService
}

func Handler(otpService otp.OTPService, authService AuthService) interfaces.HandlerInterface {
	return &authHandler{
		otpService:  otpService,
		authService: authService,
	}
}

// RegisterRoutes registers the routes for the auth handler
//
//   - router: the router group to register the routes to
//   - middlewares: the middlewares to use for the routes
//
// Routes:
//
//   - POST /login - handleAuth
//   - POST /verify - handleVerifyAuth
func (a *authHandler) RegisterRoutes(authRoute *gin.RouterGroup, middlewares []gin.HandlerFunc) {
	authRoute.Use(middlewares...)
	authRoute.POST("/", a.handleAuth)
	authRoute.POST("/verify", a.handleVerifyAuth)
}

func (a *authHandler) handleAuth(context *gin.Context) {
	var user User
	if err := context.BindJSON(&user); err != nil {
		response.BadRequest(context, err.Error())
		return
	}

	otp, err := a.otpService.RequestOTP(user.Email, user.Name)
	if err != nil {
		response.InternalServerError(context)
		return
	}

	response.Success(context, "Please check your email for the OTP.", otp.ToJSON())

}

// TODO: delete otp after verification
func (a *authHandler) handleVerifyAuth(context *gin.Context) {
	var _otp otp.Otp
	if err := context.BindJSON(&_otp); err != nil {
		response.BadRequest(context, err.Error())
		return
	}

	_otp, err := a.otpService.VerifyOTP(_otp.Email, _otp.Code, _otp.RequestId)
	if err != nil {
		response.InternalServerError(context)
		return
	}

	isVerified := _otp != (otp.Otp{})

	if !isVerified || _otp.IsExpired() {
		response.DefaultResponse(context, http.StatusNotFound, "Invalid OTP", nil)
		return
	}

	newUser := New(_otp.Name, _otp.Email, true)
	err = a.authService.CreateUser(*newUser)
	if err != nil {
		response.InternalServerError(context)
		return
	}

	token, err := a.authService.GenerateToken(*newUser)
	if err != nil {
		response.InternalServerError(context)
		return
	}

	response.Success(context, "Authenticated", gin.H{
		"token": token,
	})

}
