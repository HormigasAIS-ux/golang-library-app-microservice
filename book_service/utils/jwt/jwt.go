package jwt_util

import (
	"book_service/domain/dto"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWT(tokenString string, secretKey string) (*dto.CurrentUser, error) {
	var JWT_SIGNATURE_KEY = []byte(secretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}

		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	sub, _ := claims["sub"].(string)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	email, _ := claims["email"].(string)

	return &dto.CurrentUser{
		UUID:     sub,
		Email:    email,
		Username: username,
		Role:     role,
	}, nil
}
