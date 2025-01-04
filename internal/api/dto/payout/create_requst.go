package dto

type PayoutRequest struct {
	AccountName   string ` binding:"required,gte=3"`
	AccountNumber string ` binding:"required"`
	BankName      string ` binding:"required,gte=3"`
	BankCode      string ` binding:"required"`
}
