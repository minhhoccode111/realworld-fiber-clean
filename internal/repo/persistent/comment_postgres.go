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

// CommentRepo implements comment persistence against Postgres.
type CommentRepo struct {
	*postgres.Postgres
}

// NewCommentRepo constructs a new CommentRepo.
func NewCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

// StoreCreate inserts a comment for the given article slug and returns its ID.
func (r *CommentRepo) StoreCreate(
	ctx context.Context,
	slug string,
	dto *entity.Comment,
) (string, error) {
	sql, args, err := r.Builder.
		Insert("comments").
		Columns("author_id", "article_id", "body").
		Values(dto.AuthorID, squirrel.Expr("(select id from articles where slug = ?)", slug), dto.Body).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("CommentRepo - StoreCreate - r.Builder: %w", err)
	}

	var id string

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("CommentRepo - StoreCreate - r.Pool.QueryRow: %w", entity.ErrNoRows)
	}

	if err != nil {
		return "", fmt.Errorf("CommentRepo - StoreCreate - r.Pool.QueryRow: %w", err)
	}

	return id, nil
}

// GetDetailByID returns detailed information for a specific comment.
func (r *CommentRepo) GetDetailByID(
	ctx context.Context,
	userID, commentID string,
) (*entity.CommentDetail, error) {
	query := `
		select c.id, c.body, c.created_at, c.updated_at,
		  u.username, u.bio, u.image,
		  (select exists (
			select 1 from follows
			where follower_id::text = $1
			and following_id = c.author_id
		  )) as following,
		  (select count(distinct(follower_id)) from follows
		  where following_id = c.author_id) as followers_count
		from comments c
		left join users u on u.id = c.author_id
		left join articles a on a.id = c.article_id
		where c.deleted_at is null
		and a.deleted_at is null
		and c.id = $2;
	`
	args := []any{userID, commentID}

	/*
		example query output:
		                  id                  |  body  |          created_at           |          updated_at           | username | bio | image | following
		--------------------------------------+--------+-------------------------------+-------------------------------+----------+-----+-------+-----------
		 da1b0dc3-e2a5-4930-9e5d-1dd6f7884717 | body 0 | 2025-10-06 13:45:43.116717+00 | 2025-10-06 13:45:43.116717+00 | asd0     |     |       | f
	*/

	c := entity.CommentDetail{}

	err := r.Pool.QueryRow(ctx, query, args...).Scan(
		&c.ID,
		&c.Body,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Author.Username,
		&c.Author.Bio,
		&c.Author.Image,
		&c.Author.Following,
		&c.Author.FollowersCount,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"CommentRepo - GetDetailByID - r.Pool.QueryRow: %w",
			err,
		)
	}

	return &c, nil
}

// GetList returns a paginated list of comments and their total count.
func (r *CommentRepo) GetList(
	ctx context.Context,
	userID, slug string,
	limit, offset uint64,
) (comments []entity.CommentDetail, total uint64, err error) {
	query := `
		select c.id, c.body, c.created_at, c.updated_at,
		  u.username, u.bio, u.image,
		  (select exists (
			select 1 from follows
			where follower_id::text = $1
			and following_id = c.author_id
		  )) as following,
		  (select count(distinct(follower_id)) from follows
		  where following_id = c.author_id) as followers_count,
		  count(*) over() as comments_count
		from comments c
		left join users u on u.id = c.author_id
		left join articles a on a.id = c.article_id
		where c.deleted_at is null
		and a.deleted_at is null
		and a.slug = $2
		group by a.id, u.id, c.id
		order by c.created_at
		limit $3
		offset $4;
	`
	args := []any{userID, slug, limit, offset}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - GetList - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	comments = []entity.CommentDetail{}
	for rows.Next() {
		var c entity.CommentDetail

		err := rows.Scan(
			&c.ID,
			&c.Body,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Author.Username,
			&c.Author.Bio,
			&c.Author.Image,
			&c.Author.Following,
			&c.Author.FollowersCount,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("CommentRepo - GetList - rows.Scan: %w", err)
		}

		comments = append(comments, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - GetList - rows.Err: %w", err)
	}

	return comments, total, nil
}

// GetBasicByID returns basic comment data for authorization or deletes.
func (r *CommentRepo) GetBasicByID(ctx context.Context, commentID string) (*entity.Comment, error) {
	sql, args, err := r.Builder.
		Select("id, author_id, article_id, body, created_at, updated_at").
		From("comments").
		Where(squirrel.Eq{"id": commentID}).
		Where("deleted_at is null").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetBasicByID - r.Builder: %w", err)
	}

	c := entity.Comment{}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&c.ID,
		&c.AuthorID,
		&c.ArticleID,
		&c.Body,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"CommentRepo - GetBasicByID - r.Pool.QueryRow: %w",
			err,
		)
	}

	return &c, nil
}

// StoreDelete soft-deletes a comment bound to an article slug.
func (r *CommentRepo) StoreDelete(ctx context.Context, slug, commentID string) (err error) {
	sql, args, err := r.Builder.
		Update("comments").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Expr(`exists (
			select 1 from articles
			where id = article_id
			and slug = ?
			and deleted_at is null)`, slug)).
		Where(squirrel.Eq{"id": commentID}).
		Where("deleted_at is null").
		ToSql()
	if err != nil {
		return fmt.Errorf("CommentRepo - StoreDelete - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CommentRepo - StoreDelete - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("CommentRepo - StoreDelete - r.Pool.Exec: %w", entity.ErrNoEffect)
	}

	return nil
}
