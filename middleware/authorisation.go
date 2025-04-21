package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil
	}

	return tokenString, nil

}

func VerifyJWT(tokenSting string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenSting, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RoleMiddleware(r *http.Request) (string, error) {
	// get the authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header required")
	}

	// trim the header from bearer
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// verify the token

	token, err := VerifyJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	// parse the token against claims

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("could not parse claims")
	}

	// check the role in claims
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found in claims")
	}

	return role, nil
}

func EditorOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := RoleMiddleware(r)
		if err != nil || role == "editor" {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is empty")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := VerifyJWT(tokenString)
	if err != nil {
		return 0, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to parse claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to find userID in claims")
	}

	return uint(userID), nil
}
