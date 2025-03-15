package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/xneverov/todo-list/internal/config"
)

var secretKey = []byte("axolotl")

func generateToken(passwordHash string) (string, error) {
	claims := jwt.MapClaims{
		"passwordHash": passwordHash,
		"exp":          time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func validateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	password := config.Get("TODO_PASSWORD")

	passwordTokenHash, ok := claims["passwordHash"].(string)
	if !ok {
		return false, fmt.Errorf("password hash not found in token")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordTokenHash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
