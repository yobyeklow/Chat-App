package repository

import (
	"context"
	"web_socket/internal/database/sqlc"

	"github.com/google/uuid"
)

type SQLUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

func (ur *SQLUserRepository) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	userData, err := ur.db.CreateUser(ctx, arg)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) FindUserByEmail(ctx context.Context, userEmail string) (sqlc.User, error) {
	userData, err := ur.db.FindUserByEmail(ctx, userEmail)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) FindUserByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.FindUserByUUID(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) SoftDeleteUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.SoftDelete(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) HardDeleteUser(ctx context.Context, userUuid uuid.UUID) error {
	_, err := ur.db.HardDelete(ctx, userUuid)
	if err != nil {
		return err
	}
	return nil
}
func (ur *SQLUserRepository) RestoreUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.RestoreUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
