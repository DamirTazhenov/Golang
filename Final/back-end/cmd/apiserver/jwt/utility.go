package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"shop/internal/config"
	"time"
)

var jwtKey = []byte(config.NewConfig().JWTSecretKey)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT parses and validates the JWT token
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
