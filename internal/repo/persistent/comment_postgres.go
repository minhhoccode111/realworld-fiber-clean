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

type CommentRepo struct {
	*postgres.Postgres
}

func NewCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

func (r *CommentRepo) StoreCreate(
	ctx context.Context,
	slug string,
	dto entity.Comment,
) (string, error) {
	sql, args, err := r.Builder.
		Insert("comments").
		Columns("author_id", "article_id", "body").
		Values(dto.AuthorId, squirrel.Expr("(select id from articles where slug = ?)", slug), dto.Body).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("CommentRepo - StoreCreate - r.Builder: %w", err)
	}

	var id string
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("notfound")
		}
		return "", fmt.Errorf("CommentRepo - StoreCreate - r.Pool.QueryRow: %w", err)
	}

	return id, nil
}

func (r *CommentRepo) GetDetailById(
	ctx context.Context,
	userId, commentId string,
) (entity.CommentDetail, error) {
	query := `
		select c.id, c.body, c.created_at, c.updated_at,
		  u.username, u.bio, u.image,
		  (select exists (
			select 1 from follows
			where follower_id::text = $1
			and following_id = c.author_id
		  )) as following
		from comments c
		left join users u on u.id = c.author_id
		left join articles a on a.id = c.article_id
		where c.deleted_at is null
		and a.deleted_at is null
		and c.id = $2;
	`
	var args = []any{userId, commentId}

	/*
		example query output:
		                  id                  |  body  |          created_at           |          updated_at           | username | bio | image | following
		--------------------------------------+--------+-------------------------------+-------------------------------+----------+-----+-------+-----------
		 da1b0dc3-e2a5-4930-9e5d-1dd6f7884717 | body 0 | 2025-10-06 13:45:43.116717+00 | 2025-10-06 13:45:43.116717+00 | asd0     |     |       | f
	*/

	c := entity.CommentDetail{}
	err := r.Pool.QueryRow(ctx, query, args...).Scan(
		&c.Id,
		&c.Body,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Author.Username,
		&c.Author.Bio,
		&c.Author.Image,
		&c.Author.Following,
	)
	if err != nil {
		return entity.CommentDetail{}, fmt.Errorf(
			"CommentRepo - GetDetailById - r.Pool.QueryRow: %w",
			err,
		)
	}

	return c, nil
}
