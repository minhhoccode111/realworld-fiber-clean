package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type ProfileRepo struct {
	*postgres.Postgres
}

func NewProfileRepo(pg *postgres.Postgres) *ProfileRepo {
	return &ProfileRepo{pg}
}

func (r *ProfileRepo) IsExisted(ctx context.Context, username string) error {
	sql, args, err := r.Builder.
		Select().
		Column(squirrel.Expr("exists (select 1 from users where username = ?)", username)).
		From("users").
		ToSql()
	if err != nil {
		return fmt.Errorf("ProfileRepo - IsExisted - r.Builder: %w", err)
	}

	var isExisted bool
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&isExisted)
	if err != nil {
		return fmt.Errorf("ProfileRepo - IsExisted - r.Pool.QueryRow: %w", err)
	}
	if !isExisted {
		return fmt.Errorf("ProfileRepo - IsExisted - r.Pool.QueryRow: %w", entity.ErrNoRows)
	}

	return nil
}

func (r *ProfileRepo) GetDetail(
	ctx context.Context,
	userId, username string,
) (entity.ProfilePreview, error) {
	sql, args, err := r.Builder.
		Select(
			"username",
			"bio",
			"image",
		).
		Column(
			squirrel.Expr(
				`(select exists
					(select 1 from follows
						where follower_id::text = ?
						and following_id = (select id from users where username = ?))
				) as following`,
				userId,
				username,
			),
		).
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return entity.ProfilePreview{}, fmt.Errorf("ProfileRepo - GetDetail - r.Builder: %w", err)
	}

	var e entity.ProfilePreview
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&e.Username, &e.Bio, &e.Image, &e.Following)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.ProfilePreview{}, fmt.Errorf(
			"ProfileRepo - GetDetail - r.Pool.QueryRow: %w",
			entity.ErrNoRows,
		)
	}
	if err != nil {
		return entity.ProfilePreview{}, fmt.Errorf(
			"ProfileRepo - GetDetail - r.Pool.QueryRow: %w",
			err,
		)
	}

	return e, nil
}

func (r *ProfileRepo) StoreCreate(ctx context.Context, userId, username string) error {
	sql, args, err := r.Builder.
		Insert("follows").
		Columns("follower_id", "following_id").
		Values(userId, squirrel.Expr("(select id from users where username = ?)", username)).
		Suffix("on conflict do nothing").
		ToSql()
	if err != nil {
		return fmt.Errorf("ProfileRepo - StoreCreate - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ProfileRepo - StoreCreate - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *ProfileRepo) StoreDelete(ctx context.Context, userId, username string) error {
	sql, args, err := r.Builder.
		Delete("follows").
		Where(squirrel.Expr("follower_id = ?", userId)).
		Where(squirrel.Expr("following_id = (select id from users where username = ?)", username)).
		ToSql()
	if err != nil {
		return fmt.Errorf("ProfileRepo - StoreDelete - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ProfileRepo - StoreDelete - r.Pool.Exec: %w", err)
	}

	return nil
}
