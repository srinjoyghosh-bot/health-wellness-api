package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (uint, error)
}

type jwtService struct {
	secretKey string
	expires   time.Duration
}

func NewJWTService(secretKey string, expires time.Duration) JWTService {
	return &jwtService{
		secretKey: secretKey,
		expires:   expires,
	}
}

func (s *jwtService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.expires).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		log.Println("Error in parse", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, err
}
