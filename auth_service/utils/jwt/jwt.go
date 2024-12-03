package jwt_util

import (
	"auth_service/domain/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJwtToken(user *model.User, secretKey string, expHours int, tokenId *string) (string, error) {
	JWT_SIGNATURE_KEY := []byte(secretKey)

	claims := jwt.MapClaims{
		"sub":      user.UUID.String(),
		"username": user.Username,
		"fullname": user.Fullname,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
	}

	if tokenId != nil {
		claims["token_id"] = *tokenId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
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

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
