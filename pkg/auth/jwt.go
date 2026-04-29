package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"web_socket/internal/common/database/sqlc"
	"web_socket/internal/common/utils"
	"web_socket/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache cache.RedisService
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"user_role"`
}
type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwt_secret_key    = []byte(utils.GetEnv("JWT_SECRET_KEY", "01594057772770059127161404529220"))
	jwt_encrypted_key = []byte(utils.GetEnv("JWT_SECRET_KEY", "01594057772770059127161404523216"))
)

const (
	AccessTokenTTL  = 6 * time.Hour
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(cache cache.RedisService) TokenService {
	return &JWTService{
		cache: cache,
	}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserRole,
	}

	raw_data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Parse Raw Data: %w", err)
	}
	encrypted, err := utils.EncryptAES(raw_data, jwt_secret_key)
	if err != nil {
		return "", fmt.Errorf("Encrypt payload: %w", err)
	}

	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		"iat":  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret_key)
}
func (js *JWTService) ParseToken(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return jwt_secret_key, nil
	})
	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token claims", utils.ErrCodeUnauthorized)
	}
	return token, claims, nil
}
func (js *JWTService) DecryptAccessTokenPayload(tokenStr string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenStr)
	if err != nil {
		return nil, utils.WrapError("Can't parse JWT token", utils.ErrCodeInternal, err)
	}
	encryptedPayload, ok := claims["data"].(string)
	if !ok {
		return nil, utils.WrapError("Encoded data not found", utils.ErrCodeInternal, err)
	}

	decryptBytes, err := utils.DecryptAES(encryptedPayload, jwt_secret_key)
	if err != nil {
		return nil, utils.WrapError("Can't decode data", utils.ErrCodeInternal, err)
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(decryptBytes, &payload); err != nil {
		return nil, utils.WrapError("Invalid data format", utils.ErrCodeInternal, err)
	}

	return &payload, nil
}

func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, nil
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidJWTToken(tokenStr string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + tokenStr
	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil || refreshToken.Revoked == true || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError("Can't get refresh token", utils.ErrCodeInternal, err)
	}
	return refreshToken, nil
}
func (js *JWTService) RevokeToken(tokenStr string) error {
	cacheKey := "refresh_token:" + tokenStr
	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil || refreshToken.Revoked == true || refreshToken.ExpiresAt.Before(time.Now()) {
		return utils.WrapError("Can't get refresh token", utils.ErrCodeInternal, err)
	}
	refreshToken.Revoked = true
	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
