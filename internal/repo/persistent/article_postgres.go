package persistent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type ArticleRepo struct {
	*postgres.Postgres
}

func NewArticleRepo(pg *postgres.Postgres) *ArticleRepo {
	return &ArticleRepo{pg}
}

func (r *ArticleRepo) StoreCreate(ctx context.Context, dto entity.Article, tags []string) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - r.Pool.Begin: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r) // continue the panicking
		} else if err != nil {
			tx.Rollback(ctx)
		}
	}()

	sql, args, err := r.Builder.
		Insert("articles").
		Columns("author_id", "slug", "title", "description", "body").
		Values(dto.AuthorId, dto.Slug, dto.Title, dto.Description, dto.Body).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - r.Builder: %w", err)
	}

	// TODO: insert article and tags can be done concurrently
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&dto.Id)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - r.Pool.QueryRow: %w", err)
	}

	ids, err := r.StoreTagsList(ctx, tags)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - r.StoreTagsList: %w", err)
	}

	err = r.StoreArticleTagsList(ctx, dto.Id, ids)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - r.StoreArticleTagsList: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreCreate - tx.Commit: %w", err)
	}

	return nil
}

func (r *ArticleRepo) GetDetailBySlug(ctx context.Context, userId, slug string,
) (entity.ArticleDetail, error) {
	sql, args, err := r.Builder.
		Select(
			"a.slug",
			"a.title",
			"a.description",
			"a.body",
			"a.created_at",
			"a.updated_at",
			"coalesce(array_agg(distinct t.name) filter (where t.name is not null), '{}') as tags",
		).
		Column(squirrel.Expr("(select exists (select 1 from favorites where article_id = a.id and user_id::text = ?)) as favorited", userId)).
		Column(squirrel.Expr("(count(distinct f.user_id)) as favorites_count")).
		Columns(
			"u.username",
			"u.bio",
			"u.image",
		).
		Column(squirrel.Expr("(select exists (select 1 from follows where a.author_id = following_id and follower_id::text = ?)) as following", userId)).
		From("articles a").
		LeftJoin("users u on a.author_id = u.id").
		LeftJoin("article_tags at on at.article_id = a.id").
		LeftJoin("tags t on at.tag_id = t.id").
		LeftJoin("favorites f on f.article_id = a.id").
		Where("a.deleted_at is null").
		Where(squirrel.Eq{"a.slug": slug}).
		GroupBy("a.id", "u.id").
		ToSql()
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleRepo - GetDetailBySlug - r.Builder: %w",
			err,
		)
	}

	a := entity.ArticleDetail{}
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&a.Slug,
		&a.Title,
		&a.Description,
		&a.Body,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.TagList,
		&a.Favorited,
		&a.FavoritesCount,
		&a.Author.Username,
		&a.Author.Bio,
		&a.Author.Image,
		&a.Author.Following,
	)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleRepo - GetDetailBySlug - row.Scan: %w",
			err,
		)
	}

	return a, nil
}

func (r *ArticleRepo) StoreTagsList(ctx context.Context, tags []string,
) (ids []string, err error) {
	if len(tags) == 0 {
		return []string{}, nil
	}

	builder := r.Builder.
		Insert("tags").
		Columns("name")

	for _, tag := range tags {
		builder = builder.Values(tag)
	}

	builder = builder.Suffix("on conflict (name) do update set name=EXCLUDED.name returning id")

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ArticleRepo - StoreTagsList - builder.ToSql: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ArticleRepo - StoreTagsList - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	ids = []string{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *ArticleRepo) StoreArticleTagsList(ctx context.Context, articleId string, tagIds []string,
) error {
	if len(tagIds) == 0 {
		return nil
	}

	builder := r.Builder.
		Insert("article_tags").
		Columns("article_id", "tag_id")

	for _, tagId := range tagIds {
		builder = builder.Values(articleId, tagId)
	}

	builder = builder.Suffix("on conflict do nothing")

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreArticleTagsList - builder.ToSql: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ArticleRepo - StoreArticleTagsList - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *ArticleRepo) CanSlugBeUsed(ctx context.Context, articleId, slug string) (bool, error) {
	// query to check if slug is already existed
	query, _, err := r.Builder.
		Select("exists (select 1 from articles where id::text <> ? and slug = ?)").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("ArticleRepo - CanSlugBeUsed - r.Builder: %w", err)
	}

	var existed bool
	err = r.Pool.QueryRow(ctx, query, articleId, slug).Scan(&existed)
	if err != nil {
		return false, fmt.Errorf("ArticleRepo - CanSlugBeUsed - r.Builder: %w", err)
	}

	return !existed, nil
}
