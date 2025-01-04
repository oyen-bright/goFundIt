// internal/tests/integration/auth_test.go

package integration

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/oyen-bright/goFundIt/config/providers"
// 	"github.com/oyen-bright/goFundIt/internal/api/handlers"
// 	"github.com/oyen-bright/goFundIt/internal/api/middlewares"
// 	"github.com/oyen-bright/goFundIt/internal/models"
// 	postgress "github.com/oyen-bright/goFundIt/internal/repositories/postgres"
// 	"github.com/oyen-bright/goFundIt/internal/services"
// 	"github.com/oyen-bright/goFundIt/pkg/email"
// 	"github.com/oyen-bright/goFundIt/pkg/encryption"
// 	"github.com/oyen-bright/goFundIt/pkg/jwt"
// 	"github.com/oyen-bright/goFundIt/pkg/logger"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// type AuthIntegrationTestSuite struct {
// 	suite.Suite
// 	db     *gorm.DB
// 	router *gin.Engine
// 	jwt    jwt.Jwt
// }

// func (suite *AuthIntegrationTestSuite) SetupSuite() {
// 	gin.SetMode(gin.TestMode)

// 	// Initialize in-memory SQLite database
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	suite.Require().NoError(err)

// 	// Run migrations
// 	err = db.AutoMigrate(&models.User{}, &models.Contributor{})
// 	suite.Require().NoError(err)

// 	suite.db = db

// 	// Initialize dependencies
// 	encryptor := encryption.New([]string{"test-secret"})
// 	emailer := email.New(providers.EmailSMTP, email.EmailConfig{})
// 	jwtService := jwt.New("test-secret")
// 	logger := logger.New()

// 	// Initialize repositories
// 	authRepo := postgress.NewAuthRepository(db)
// 	otpRepo := postgress.NewOTPRepository(db)

// 	// Initialize services
// 	otpService := services.NewOTPService(otpRepo, emailer, *encryptor, logger) // You might want to mock this for deterministic testing
// 	authService := services.NewAuthService(authRepo, otpService, *encryptor, jwtService, logger)

// 	// Initialize handlers
// 	authHandler := handlers.NewAuthHandler(authService)

// 	// Setup router with middleware
// 	router := gin.New()
// 	// router.Use(middlewares.Cors())

// 	// Define routes
// 	router.POST("/auth", authHandler.HandleAuth)
// 	router.POST("/auth/verify", authHandler.HandleVerifyAuth)

// 	suite.router = router
// 	suite.jwt = jwtService
// }

// func (suite *AuthIntegrationTestSuite) TearDownSuite() {
// 	// Cleanup database
// 	sqlDB, err := suite.db.DB()
// 	suite.Require().NoError(err)
// 	sqlDB.Close()
// }

// func (suite *AuthIntegrationTestSuite) TestAuthFlow() {
// 	// Test cases for the complete auth flow
// 	testCases := []struct {
// 		name          string
// 		authPayload   map[string]interface{}
// 		verifyPayload map[string]interface{}
// 		expectedCodes []int
// 		checkResponse func(*testing.T, *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "successful auth flow",
// 			authPayload: map[string]interface{}{
// 				"email": "test@example.com",
// 				"name":  "Test User",
// 			},
// 			verifyPayload: map[string]interface{}{
// 				"email":      "test@example.com",
// 				"code":       "123456", // Note: In real tests, this should match what OTP service generates
// 				"request_id": "",       // Will be filled from auth response
// 			},
// 			expectedCodes: []int{http.StatusOK, http.StatusOK},
// 			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
// 				var response map[string]interface{}
// 				err := json.NewDecoder(w.Body).Decode(&response)
// 				assert.NoError(t, err)
// 				assert.Contains(t, response, "data")
// 				data := response["data"].(map[string]interface{})
// 				assert.Contains(t, data, "token")
// 			},
// 		},
// 		{
// 			name: "invalid email format",
// 			authPayload: map[string]interface{}{
// 				"email": "invalid-email",
// 				"name":  "Test User",
// 			},
// 			verifyPayload: nil,
// 			expectedCodes: []int{http.StatusBadRequest},
// 			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
// 				var response map[string]interface{}
// 				err := json.NewDecoder(w.Body).Decode(&response)
// 				assert.NoError(t, err)
// 				assert.Contains(t, response["message"], "Invalid inputs")
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.T().Run(tc.name, func(t *testing.T) {
// 			// Step 1: Request Auth
// 			if tc.authPayload != nil {
// 				authJSON, err := json.Marshal(tc.authPayload)
// 				assert.NoError(t, err)

// 				w := httptest.NewRecorder()
// 				req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(authJSON))
// 				req.Header.Set("Content-Type", "application/json")
// 				suite.router.ServeHTTP(w, req)

// 				assert.Equal(t, tc.expectedCodes[0], w.Code)

// 				if w.Code == http.StatusOK && tc.verifyPayload != nil {
// 					// Extract request_id from response
// 					var response map[string]interface{}
// 					err := json.NewDecoder(w.Body).Decode(&response)
// 					assert.NoError(t, err)
// 					data := response["data"].(map[string]interface{})
// 					tc.verifyPayload["request_id"] = data["request_id"]

// 					// Step 2: Verify Auth
// 					verifyJSON, err := json.Marshal(tc.verifyPayload)
// 					assert.NoError(t, err)

// 					w = httptest.NewRecorder()
// 					req = httptest.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(verifyJSON))
// 					req.Header.Set("Content-Type", "application/json")
// 					suite.router.ServeHTTP(w, req)

// 					assert.Equal(t, tc.expectedCodes[1], w.Code)
// 					tc.checkResponse(t, w)
// 				}
// 			}
// 		})
// 	}
// }

// // Additional test for protected routes
// func (suite *AuthIntegrationTestSuite) TestProtectedRoutes() {
// 	// Setup protected route
// 	suite.router.GET("/protected", middlewares.Auth(suite.jwt), func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "success"})
// 	})

// 	tests := []struct {
// 		name         string
// 		setupAuth    func() string
// 		expectedCode int
// 	}{
// 		{
// 			name: "access with valid token",
// 			setupAuth: func() string {
// 				token, _ := suite.jwt.GenerateToken(1, "test@example.com", "test-handle")
// 				return token
// 			},
// 			expectedCode: http.StatusOK,
// 		},
// 		{
// 			name: "access without token",
// 			setupAuth: func() string {
// 				return ""
// 			},
// 			expectedCode: http.StatusUnauthorized,
// 		},
// 	}

// 	for _, tt := range tests {
// 		suite.T().Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodGet, "/protected", nil)

// 			if token := tt.setupAuth(); token != "" {
// 				req.Header.Set("Authorization", "Bearer "+token)
// 			}

// 			suite.router.ServeHTTP(w, req)
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 		})
// 	}
// }

// func TestAuthIntegration(t *testing.T) {
// 	suite.Run(t, new(AuthIntegrationTestSuite))
// }

// // func setupTestDB() (*gorm.DB, error) {
// /*
//    dsn := fmt.Sprintf(
//        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
//        os.Getenv("TEST_DB_HOST"),
//        os.Getenv("TEST_DB_USER"),
//        os.Getenv("TEST_DB_PASSWORD"),
//        os.Getenv("TEST_DB_NAME"),
//        os.Getenv("TEST_DB_PORT"),
//    )

//    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//    if err != nil {
//        return nil, err
//    }

//    // Clean database
//    db.Exec("TRUNCATE users, contributions RESTART IDENTITY CASCADE")

//    // Run migrations
//    err = db.AutoMigrate(&models.User{}, &models.Contribution{})
//    if err != nil {
//        return nil, err
//    }

//    return db, nil
// */
// // }
