package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey      string
	defaultIssuer  string
	expirationTime int64
}

func NewJwtTokenService(secretKey, defaultIssuer string, expirationTime int64) *JwtService {
	return &JwtService{
		secretKey:      secretKey,
		defaultIssuer:  defaultIssuer,
		expirationTime: expirationTime,
	}
}

func (s *JwtService) SignToken(subject string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    s.defaultIssuer,
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (s *JwtService) ValidateToken(tokenString string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return &JwtToken{
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
			ExpiresAt: claims.ExpiresAt.Time.Unix(),
		}, nil
	} else {
		return nil, ErrInvalidToken
	}
}
