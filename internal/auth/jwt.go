package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
			Subject: userID.String(),
		},
	)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("not a valid token")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims) ; ok {
		parsedId, pErr := uuid.Parse(claims.Subject)
		if pErr != nil {
			return uuid.Nil, pErr
		}
		return parsedId, nil
	} else {
		return uuid.Nil, errors.New("error parsing this JWT")
	}
}