package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

type TokenType string

const (
	TokenTypeAccess  TokenType = "chirpy-access"
	TokenTypeRefresh TokenType = "chirpy-refresh"
)

func MakeJWT(userID int, tokenSecret string, expiresIn time.Duration, tokenType TokenType) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	tokenIssuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if tokenIssuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}

func RefreshToken(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	tokenIssuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if tokenIssuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid issuer")
	}

	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		return "", err
	}

	newToken, err := MakeJWT(userIDInt, tokenSecret, time.Hour, TokenTypeAccess)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splithAuth := strings.Split(authHeader, " ")
	if len(splithAuth) < 2 || splithAuth[0] != "Bearer" {
		return "", errors.New("malformed authoriztion header")
	}

	return splithAuth[1], nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splithAuth := strings.Split(authHeader, " ")
	if len(splithAuth) < 2 || splithAuth[0] != "ApiKey" {
		return "", errors.New("malformed authoriztion header")
	}

	return splithAuth[1], nil
}
