package postgress

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsSave(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewAnalyticsRepository(db)

	now := time.Now()
	analytics := &models.PlatformAnalytics{
		FiatStats:      map[string]models.CurrencyStats{"USD": {Amount: 100}},
		CryptoStats:    map[string]models.CurrencyStats{"BTC": {Amount: 1}},
		PaymentMethods: models.PaymentMethodStats{Fiat: 5, Manual: 3},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := repo.Save(analytics)
	assert.NoError(t, err)

	// Verify save
	saved, err := repo.Get(now)
	assert.NoError(t, err)
	assert.Equal(t, analytics.FiatStats["USD"].Amount, saved.FiatStats["USD"].Amount)
	assert.Equal(t, analytics.CryptoStats["BTC"].Amount, saved.CryptoStats["BTC"].Amount)
	assert.Equal(t, analytics.PaymentMethods.Fiat, saved.PaymentMethods.Fiat)
	assert.Equal(t, analytics.PaymentMethods.Manual, saved.PaymentMethods.Manual)
}

func TestAnalyticsGet(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewAnalyticsRepository(db)

	tests := []struct {
		name    string
		setup   func() time.Time
		wantErr bool
	}{
		{
			name: "get existing analytics",
			setup: func() time.Time {
				now := time.Now()
				analytics := &models.PlatformAnalytics{
					FiatStats:      map[string]models.CurrencyStats{"EUR": {Amount: 200}},
					CryptoStats:    map[string]models.CurrencyStats{"ETH": {Amount: 2}},
					PaymentMethods: models.PaymentMethodStats{Fiat: 10},
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				require.NoError(t, repo.Save(analytics))
				return now
			},
			wantErr: false,
		},
		{
			name: "get non-existing date creates new",
			setup: func() time.Time {
				return time.Now().AddDate(0, 0, 1) // tomorrow
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := tt.setup()
			analytics, err := repo.Get(date)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, analytics)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, analytics)
				assert.Equal(t, date.Year(), analytics.CreatedAt.Year())
				assert.Equal(t, date.Month(), analytics.CreatedAt.Month())
				assert.Equal(t, date.Day(), analytics.CreatedAt.Day())
			}
		})
	}
}
