package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOTPAdd(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewOTPRepository(db)

	tests := []struct {
		name    string
		otp     *models.Otp
		wantErr bool
	}{
		{
			name:    "valid otp",
			otp:     models.NewOTP("test@example.com"),
			wantErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			err := repo.Add(tt.otp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestOTPGetByEmailAndReference(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewOTPRepository(db)

	testOtp := models.NewOTP("test@example.com")

	err := repo.Add(testOtp)
	require.NoError(t, err)

	requestId := testOtp.RequestId
	code := testOtp.Code

	tests := []struct {
		name      string
		email     string
		requestId string
		code      string
		want      bool
		wantErr   bool
	}{
		{
			name:      "valid otp",
			email:     testOtp.Email,
			code:      code,
			requestId: requestId,
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Invalid request_id",
			email:     testOtp.Email,
			code:      code,
			requestId: "invalid-request-id",
			want:      true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := repo.GetByEmailAndReference(tt.email, tt.requestId)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, otp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, otp)

			}
		})
	}
}

func TestOTPDelete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewOTPRepository(db)

	testOtp := models.NewOTP("test@example.com")

	err := repo.Add(testOtp)
	require.NoError(t, err)

	tests := []struct {
		name    string
		otp     *models.Otp
		wantErr bool
	}{
		{
			name:    "valid otp",
			otp:     testOtp,
			wantErr: false,
		},
		{
			name:    "invalid otp",
			otp:     models.NewOTP("testUser@example.com"),
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.otp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
func TestOTPInvalidateOtherOTPs(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewOTPRepository(db)

	// Create multiple OTPs for same email
	email := "test@example.com"
	otp1 := models.NewOTP(email)
	otp2 := models.NewOTP(email)
	otp3 := models.NewOTP(email)

	require.NoError(t, repo.Add(otp1))
	require.NoError(t, repo.Add(otp2))
	require.NoError(t, repo.Add(otp3))

	tests := []struct {
		name      string
		email     string
		code      string
		requestId string
		wantErr   bool
	}{
		{
			name:      "valid invalidation",
			email:     email,
			code:      otp1.Code,
			requestId: otp1.RequestId,
			wantErr:   false,
		},
		{
			name:      "non-existent email",
			email:     "wrong@example.com",
			code:      otp1.Code,
			requestId: otp1.RequestId,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.InvalidateOtherOTPs(tt.email, tt.code, tt.requestId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify only matching OTP remains
				var remainingOTPs []models.Otp
				err := db.Where("email = ?", tt.email).Find(&remainingOTPs).Error
				assert.NoError(t, err)

				if tt.email == email {
					assert.Equal(t, 1, len(remainingOTPs))
					assert.Equal(t, tt.code, remainingOTPs[0].Code)
					assert.Equal(t, tt.requestId, remainingOTPs[0].RequestId)
				} else {
					assert.Equal(t, 0, len(remainingOTPs))
				}
			}
		})
	}
}
