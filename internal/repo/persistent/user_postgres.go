package persistent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) StoreRegister(ctx context.Context, user entity.User) (entity.User, error) {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("email, username, password, bio, image").
		Values(user.Email, user.Username, user.Password, user.Bio, user.Image).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - StoreRegisterUser - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&user.Id)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - StoreRegisterUser - row.Scan: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, password, bio, image").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetUserByEmail - r.Builder: %w", err)
	}

	var user entity.User
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Bio,
		&user.Image,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetUserByEmail - row.Scan: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, userId string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, password, bio, image").
		From("users").
		Where(squirrel.Eq{"id": userId}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetUserById - r.Builder: %w", err)
	}

	var user entity.User
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Bio,
		&user.Image,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetUserById - row.Scan: %w", err)
	}

	return user, nil
}
