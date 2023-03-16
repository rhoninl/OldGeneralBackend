package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"ExpiredAt": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ValidateToken(token string) (string, bool) {
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return "", false
	}

	if claims, ok := tokenString.Claims.(jwt.MapClaims); ok &&
		tokenString.Valid &&
		claims["ExpiredAt"].(int64) > time.Now().Unix() {
		return claims["userId"].(string), true
	}

	return "", false
}
