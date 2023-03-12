package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

type JwtService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type HMacJwtService struct {
	SecretKey string
	Issuer    string
}

type QQClaims struct {
	jwt.StandardClaims
	UserID string `json:"userID"`
}

func (service *HMacJwtService) GenerateToken(userID string) (string, error) {
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

func (service *HMacJwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			logrus.Error("Invalid token")
			return nil, errors.New("invalid token")
		}
		return []byte(service.SecretKey), nil
	})
}
