package dto

type InitializePaymentResponse struct {
	Reference     string `json:"reference" binding:"required" example:"PAY-REF-123456789"`
	PaymentLink   string `json:"payment_link" binding:"required,url" example:"https://payment-gateway.com/pay/PAY-REF-123456789"`
	PaymentStatus string `json:"payment_status" binding:"required,oneof=pending completed failed" example:"pending"`
}
