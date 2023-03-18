package helper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

const (
	CONTEXT_USER_TOKEN_AUTHORIZATION_STR = "Authorization"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"ExpiredAt": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(token string) (string, bool) {
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return "", false
	}

	if claims, ok := tokenString.Claims.(jwt.MapClaims); ok &&
		tokenString.Valid &&
		claims["ExpiredAt"].(float64) > float64(time.Now().Unix()) {
		return claims["userId"].(string), true
	}

	return "", false
}

func GetUserIdFromContext(ctx context.Context) (string, error) {
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return "", errors.New("cannot get incomming context")
	}

	tokens := md.Get(CONTEXT_USER_TOKEN_AUTHORIZATION_STR)
	if len(tokens) < 1 {
		return "", errors.New("cannot get authorization")
	}

	token := tokens[0]
	userId, legal := ValidateToken(token)
	if !legal {
		return "", errors.New("token is not valid")
	}

	return userId, nil
}
