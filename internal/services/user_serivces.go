package services

import (
	"errors"
	"web_socket/internal/database/sqlc"
	"web_socket/internal/repository"
	"web_socket/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userServices struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServices{
		repo: repo,
	}
}
func (us *userServices) CreateUser(ctx *gin.Context, userInput sqlc.CreateUserParams) (sqlc.User, error) {
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
func (us *userServices) FindUserByEmail(ctx *gin.Context, userEmail string) (sqlc.User, error) {
	userData, err := us.repo.FindUserByEmail(ctx, userEmail)
	if err != nil {
		return sqlc.User{}, utils.WrapError("User is not existed", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
