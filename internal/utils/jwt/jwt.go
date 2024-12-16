package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt interface {
	GenerateToken(userID uint, email, handle string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtCfg struct {
	jwtSecret []byte
}

func New(jwtSecret string) Jwt {
	return &jwtCfg{jwtSecret: []byte(jwtSecret)}
}

type Claims struct {
	Id     uint   `json:"id"`
	Email  string `json:"email"`
	Handle string `json:"handle"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with an expiration time.
func (j jwtCfg) GenerateToken(id uint, email, handle string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires after 24 hours
	claims := &Claims{
		Id:     id,
		Email:  email,
		Handle: handle,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}

// ValidateToken validates the JWT token and returns the claims if valid.
func (j jwtCfg) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
