package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) StoreRegister(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("email, username, password").
		Values(user.Email, user.Username, user.Password).
		Suffix("returning id, image, bio, role").
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - StoreRegister - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&user.ID, &user.Image, &user.Bio, &user.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("UserRepo - StoreRegister - row.Scan: %w", entity.ErrConflict)
			}
		}

		return fmt.Errorf("UserRepo - StoreRegister - row.Scan: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, password, bio, image, role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByEmail - r.Builder: %w", err)
	}

	var user entity.User

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Bio,
		&user.Image,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UserRepo - GetUserByEmail - row.Scan: %w", entity.ErrNoRows)
		}

		return nil, fmt.Errorf("UserRepo - GetUserByEmail - row.Scan: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, password, bio, image, role").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByID - r.Builder: %w", err)
	}

	var user entity.User

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Bio,
		&user.Image,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByID - row.Scan: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) StoreUpdate(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.
		Update("users").
		Set("email", user.Email).
		Set("username", user.Username).
		Set("password", user.Password).
		Set("bio", user.Bio).
		Set("image", user.Image).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - StoreUpdate - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("UserRepo - StoreUpdate - r.Pool.Exec: %w", entity.ErrConflict)
			}
		}

		return fmt.Errorf("UserRepo - StoreUpdate - r.Pool.Exec: %w", err)
	}

	return nil
}
