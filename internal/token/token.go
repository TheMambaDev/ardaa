package token

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(secret) == 0 {
		slog.Error("JWT_SECRET not set")
	}

	InitDB()
}

func NewToken(uuid string, expire time.Time, ip net.IP) (string, error) {
	// claims ama data-da aan encode gareyn rabno
	claims := jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{
			Time: expire,
		},
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ID:      uuid,
		Subject: uuid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		slog.Error("Error creating token:", "Err", err)
		return "", err
	}

	dbToken := Token{
		ID:        claims.ID,
		Token:     tokenString,
		ExpiresAt: expire,
		Ip:        ip,
	}

	err = AddToken(dbToken)
	if err != nil {
		return "", err
	}

	// store the token in sqlite database

	return tokenString, nil
}

func Parse(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if token.Valid {
		return token.Claims.GetSubject()
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		slog.Error("That's not even a token", "Err", err)
		return "", err
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		slog.Error("Invalid signature", "Err", err)
		return "", err
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		slog.Error("Invalid Token", "Err", err)
		return "", err
	} else {
		slog.Error("Couldn't handle this token:", "Err", err)
		return "", err
	}
}
