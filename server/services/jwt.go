package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type IJwtService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JwtService struct {
	SecretKey string
	Issuer    string
}

type QQClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userID"`
}

func (service *JwtService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := &QQClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(service.SecretKey))
}

func (service *JwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, errors.New("invalid token")
		}
		return []byte(service.SecretKey), nil
	})

}
