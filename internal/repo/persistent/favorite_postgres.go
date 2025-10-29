package persistent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type FavoriteRepo struct {
	*postgres.Postgres
}

func NewFavoriteRepo(pg *postgres.Postgres) *FavoriteRepo {
	return &FavoriteRepo{pg}
}

func (r *FavoriteRepo) StoreCreate(ctx context.Context, userId, slug string) error {
	sql, args, err := r.Builder.
		Insert("favorites").
		Columns("user_id", "article_id").
		Values(squirrel.Expr(
			`?, (select id from articles where slug = ? and deleted_at is null)`,
			userId, slug,
		)).
		Suffix("on conflict do nothing").
		ToSql()
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreCreate - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreCreate - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *FavoriteRepo) StoreDelete(ctx context.Context, userId, slug string) error {
	sql, args, err := r.Builder.
		Delete("favorites").
		Where(squirrel.Expr("article_id = (select id from articles where slug = ? and deleted_at is null)", slug)).
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreDelete - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FavoriteRepo - StoreDelete - r.Pool.Exec: %w", err)
	}

	return nil
}
