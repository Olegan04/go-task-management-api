package utils

import (
	"fmt"
	"time"

	castomErr "task-manager/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID, jwtSecret []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     now.Add(24 * time.Hour).Unix(),
		"iat":     now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("jwt.GenerateToken: generate token: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, jwtSecret []byte) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt.ValidateToken: unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("jwt.ValidateToken: parse token: %w", err)
	}
	if !token.Valid {
		return uuid.Nil, castomErr.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, castomErr.ErrInvalidClaimType
	}
	userIDStr, ok := claims[string(UserIDKey)].(string)
	if !ok {
		return uuid.Nil, castomErr.ErrWithUserID
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("jwt.ValidateToken: invalid user_id format: %w", err)
	}

	return userID, nil
}
