package jwt

type JwtToken struct {
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func NewJwtToken(subject, issuer string, expiresAt, issuedAt int64) *JwtToken {
	return &JwtToken{
		Subject:   subject,
		Issuer:    issuer,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}
}
