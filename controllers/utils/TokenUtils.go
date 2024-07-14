package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(username string, apiSecret string) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN_MINUTES"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user-access"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(apiSecret))
}

func GenerateRefreshToken(username string, apiSecret string) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFESPAN_MINUTES"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["refresh-token-access"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(apiSecret))
}

func RefreshTokenValid(c *gin.Context, apiSecret string) (string, error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshTokenAccess, ok := claims["refresh-token-access"]
		if ok && refreshTokenAccess.(bool) {
			if username, ok := claims["username"]; ok {
				return username.(string), nil
			}
		}
	}
	return "", errors.New("refresh token authentication failed")
}

func ExtractUsername(c *gin.Context, apiSecret string) (string, error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userAccess, ok := claims["user-access"]
		if ok && userAccess.(bool) {
			if username, ok := claims["username"]; ok {
				return username.(string), nil
			}
		}
	}
	return "", errors.New("token authentication failed")
}

func extractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
