package dto

type VerifyAccountRequest struct {
	AccountNumber string `json:"accountNumber" binding:"required"`
	BankCode      string `json:"bankCode" binding:"required"`
}
