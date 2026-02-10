package authjwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Service {
	return &Service{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *Service) GenerateAccessToken(userID int64) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(s.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

var ErrInvalidToken = errors.New("invalid token")

func (s *Service) ParseAccessToken(tokenStr string) (int64, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// защита от "alg=none" и подмены алгоритма
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !tok.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, ErrInvalidToken
	}

	// jwt/v5 даёт числа как float64
	f, ok := sub.(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int64(f), nil
}

func (s *Service) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(7 * 24 * time.Hour).Unix(), // временно хардкод на 7 дней
		"typ": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ParseRefreshToken(tokenStr string) (int64, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !tok.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	typ, ok := claims["typ"].(string)
	if !ok || typ != "refresh" {
		return 0, ErrInvalidToken
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, ErrInvalidToken
	}

	f, ok := sub.(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int64(f), nil
}


