package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/interfaces"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/utils/response"
)

type authHandler struct {
	service AuthService
}

func Handler(service AuthService) interfaces.HandlerInterface {
	return &authHandler{
		service: service,
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

	//Validate Request
	//Email- required, email
	//Name- optional, string
	if err := context.BindJSON(&user); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", response.ExtractValidationErrors(err))
		return
	}

	otp, err := a.service.RequestAuth(user.Email, user.Name)
	if err != nil {
		response.FromError(context, err)
		return
	}

	response.Success(context, "Please check your email for the OTP.", otp.ToJSON())

}

func (a *authHandler) handleVerifyAuth(context *gin.Context) {
	var _otp otp.Otp

	//Validate Request
	//Email- required, email
	//Code- required, string
	//RequestId- required, string
	if err := context.BindJSON(&_otp); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", response.ExtractValidationErrors(err))
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
