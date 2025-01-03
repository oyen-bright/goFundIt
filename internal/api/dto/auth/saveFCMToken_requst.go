package dto

// FCMTokenRequest represents the FCM token update request
// @Description FCM token update request payload
type FCMUpdateRequest struct {
	FCMToken string `binding:"required"  example:"your-fcm-token-here"`
}
