package dto

type VerifyAuthRequest struct {
	// @Description OTP code
	// @example 123456
	Code string `json:"code" binding:"required"`
	// @Description User's email
	// @example user@example.com
	Email     string `json:"email" binding:"required,email"`
	RequestId string `json:"requestId"`
}
