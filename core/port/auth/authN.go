package port

import "github.com/golang-jwt/jwt/v5"

type JwtAuthNService interface {
	CreateToken(email string) (string, error)
	VerifyToken(token string) (jwt.MapClaims, error)
}
