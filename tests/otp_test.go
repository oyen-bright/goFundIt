package tests

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/otp/model"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&model.Otp{})
	return db
}

func TestRequestOTP(t *testing.T) {
	db := setupTestDB()
	emailer := email.NewMockEmailer()
	encryptor := encryption.New([]string{""})
	service := otp.OtpService{DB: db, Emailer: emailer, Encryptor: *encryptor}

	err := service.RequestOTP("test@example.com", "Test User")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var otp model.Otp
	db.Where("email = ?", "test@example.com").First(&otp)
	if otp.Email != "test@example.com" {
		t.Fatalf("Expected OTP email to be 'test@example.com', got %v", otp.Email)
	}
}

func TestVerifyOTP(t *testing.T) {
	db := setupTestDB()
	emailer := email.NewMockEmailer()
	encryptor := encryption.New([]string{""})
	service := otp.OtpService{DB: db, Emailer: emailer, Encryptor: *encryptor}

	otp := model.Otp{
		Email:     "test@example.com",
		Code:      "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	db.Create(&otp)

	valid, err := service.VerifyOTP("test@example.com", "123456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Fatalf("Expected OTP to be valid")
	}

	invalid, err := service.VerifyOTP("test@example.com", "654321")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if invalid {
		t.Fatalf("Expected OTP to be invalid")
	}
}
