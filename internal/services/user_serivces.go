package services

import (
	"database/sql"
	"errors"
	"web_socket/internal/database/sqlc"
	"web_socket/internal/repository"
	"web_socket/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userServices struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServices{
		repo: repo,
	}
}

func (us *userServices) FindUserByEmail(ctx *gin.Context, userEmail string) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.repo.FindUserByEmail(context, userEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to get user by email", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (us *userServices) FindUserByUUID(ctx *gin.Context, userUUID uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.repo.FindUserByUUID(context, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to get user by email", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (us *userServices) SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.repo.SoftDeleteUser(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to delete user", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (us *userServices) HardDeleteUser(ctx *gin.Context, userUuid uuid.UUID) error {
	context := ctx.Request.Context()
	err := us.repo.HardDeleteUser(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return utils.WrapError("Failed to delete user", utils.ErrCodeInternal, err)
	}
	return nil
}
func (us *userServices) RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.repo.RestoreUser(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to restore user", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
