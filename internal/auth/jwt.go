package auth

import (
	"time"
	"errors"
	"strings"
	"net/http"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := &jwt.RegisteredClaims{
    	ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
    	Issuer:    "chirpy-access",
    	IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
    	Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
	    return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			// garante que está usando HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}


func GetBearerToken(headers http.Header) (string, error) {
	
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid Authorization header format")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header must start with Bearer")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("missing token")
	}

	return token, nil
}
