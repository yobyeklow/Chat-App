package repository

import (
	"context"
	"web_socket/internal/database/sqlc"
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
