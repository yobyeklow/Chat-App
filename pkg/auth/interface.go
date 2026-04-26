package auth

import (
	"web_socket/internal/database/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	RevokeToken(tokenStr string) error
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	ValidJWTToken(tokenStr string) (RefreshToken, error)
	StoreRefreshToken(token RefreshToken) error
	ParseToken(tokenStr string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessTokenPayload(tokenStr string) (*EncryptedPayload, error)
}
