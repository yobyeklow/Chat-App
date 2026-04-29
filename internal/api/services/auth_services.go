package services

import (
	"errors"
	"strings"
	"time"
	"web_socket/internal/api/repository"
	"web_socket/internal/common/database/sqlc"
	"web_socket/internal/common/utils"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type authServices struct {
	repo         repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisService
}

func NewAuthServices(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisService) AuthService {
	return &authServices{
		repo:         repo,
		tokenService: tokenService,
		cacheService: cache,
	}
}
func (us *authServices) Register(ctx *gin.Context, userInput sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	userInput.UserEmail = utils.NormalizeString(userInput.UserEmail)

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(userInput.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError("Failed to hash password", utils.ErrCodeInternal, err)
	}
	userInput.UserPassword = string(hashed_password)
	userData, err := us.repo.CreateUser(context, userInput)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("Email already existed", utils.ErrCodeConflict)
		}
		return sqlc.User{}, utils.WrapError("Failed to create new user", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (as *authServices) Login(ctx *gin.Context, email string, password string) (string, string, int, error) {
	context := ctx.Request.Context()
	email = utils.NormalizeString(email)
	userData, err := as.repo.FindUserByEmail(context, email)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid  email or password", utils.ErrCodeUnauthorized)
	}
	//Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.UserPassword), []byte(password)); err != nil {
		return "", "", 0, utils.NewError("Invalid  email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate access token", utils.ErrCodeUnauthorized)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate refresh token", utils.ErrCodeUnauthorized)
	}
	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authServices) Logout(ctx *gin.Context, refreshTokenStr string) error {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing authorization header", utils.ErrCodeUnauthorized)
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid Access Token", utils.ErrCodeUnauthorized)
	}
	if jti, ok := claims["jti"].(string); ok {
		endUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(endUnix), 0)
		redis_key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cacheService.Set(redis_key, "revoked", ttl)
	}

	_, err = as.tokenService.ValidJWTToken(refreshTokenStr)
	if err != nil {
		return utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	err = as.tokenService.RevokeToken(refreshTokenStr)
	if err != nil {
		return utils.NewError("Unable to revoke refresh token", utils.ErrCodeInternal)
	}
	return nil
}
func (as *authServices) RefreshToken(ctx *gin.Context, refreshTokenStr string) (string, string, int, error) {
	context := ctx.Request.Context()
	token, err := as.tokenService.ValidJWTToken(refreshTokenStr)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	user_uuid, _ := uuid.Parse(token.UserUUID)

	userData, err := as.repo.FindUserByUUID(context, user_uuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate access token", utils.ErrCodeUnauthorized)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate refresh token", utils.ErrCodeUnauthorized)
	}

	if err := as.tokenService.RevokeToken(refreshTokenStr); err != nil {
		return "", "", 0, utils.NewError("Unable to revoke refresh token", utils.ErrCodeInternal)
	}
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.NewError("Can't save refresh token", utils.ErrCodeUnauthorized)
	}
	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
