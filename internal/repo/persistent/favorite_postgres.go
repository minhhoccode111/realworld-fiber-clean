package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

// FavoriteRepo implements favorite persistence backed by Postgres.
type FavoriteRepo struct {
	*postgres.Postgres
}

// NewFavoriteRepo constructs a new FavoriteRepo.
func NewFavoriteRepo(pg *postgres.Postgres) *FavoriteRepo {
	return &FavoriteRepo{pg}
}

// StoreCreate inserts a favorite link between a user and article.
func (r *FavoriteRepo) StoreCreate(ctx context.Context, userID, slug string) error {
	sql, args, err := r.Builder.
		Insert("favorites").
		Columns("user_id", "article_id").
		Values(squirrel.Expr(
			`?, (select id from articles where slug = ? and deleted_at is null)`,
			userID, slug,
		)).
		Suffix("on conflict do nothing").
		ToSql()
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreCreate - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23502" {
			return fmt.Errorf("FavoriteRepo - StoreCreate - r.Pool.Exec: %w", entity.ErrNoRows)
		}

		return fmt.Errorf("FavoriteRepo - StoreCreate - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"FavoriteRepo - StoreCreate - r.Pool.Exec: %w",
			entity.ErrNoEffect,
		)
	}

	return nil
}

// StoreDelete removes a favorite link between a user and article.
func (r *FavoriteRepo) StoreDelete(ctx context.Context, userID, slug string) error {
	sql, args, err := r.Builder.
		Delete("favorites").
		Where(squirrel.Expr("article_id = (select id from articles where slug = ? and deleted_at is null)", slug)).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreDelete - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreDelete - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"FavoriteRepo - StoreDelete - r.Pool.Exec: %w",
			entity.ErrNoEffect,
		)
	}

	return nil
}
