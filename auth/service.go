package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewJwtService() *jwtService {
	return &jwtService{}
}

var SECRET_KEY = []byte("ACESPORT_)9ajsdijoa8sdoijlasdaefw")

func (s *jwtService) GenerateToken(userId string) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenAssigned, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return tokenAssigned, err
	}

	return tokenAssigned, nil
}

func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return tokenParsed, err
	}

	return tokenParsed, nil
}
