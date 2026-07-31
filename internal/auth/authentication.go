package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassWord(password string) (string, error) {
	hash, err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash,nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	status,err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return status, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID,error) {
	token, err := jwt.ParseWithClaims(tokenString,jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret),nil
	})
	if err != nil {
		return uuid.Nil,err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil,err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil,err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("Invaid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user ID: %w",err)
	}
	
	return id,nil
}

func MakeJWT(
	userID uuid.UUID, 
	tokenSecret string, 
	expiresIn time.Duration,
) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})

	return claims.SignedString([]byte(tokenSecret))
}

func GetBearerToken(headers http.Header) (string, error){
	val := headers.Get("Authorization")
	if val == "" {
		return "", fmt.Errorf("Token doesnt exist")
	}
	res := strings.Split(val, " ")
	return res[1], nil
}

func MakeRefreshToken() string {
	hax := make([]byte, 32)
	rand.Read(hax)
	return hex.EncodeToString(hax)
}