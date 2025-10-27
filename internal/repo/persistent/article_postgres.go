package persistent

import (
	"context"
	"fmt"

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
	// NOTE: can't use squirrel because compldex queries should look complex :)
	query := `
		select a.slug, a.title, a.description, a.body, a.created_at, a.updated_at,
		  coalesce(array_agg(distinct t.name) filter (where t.name is not null), '{}') as tags,
		  (select exists
			(select 1 from favorites where article_id = a.id and user_id::text = $1)
		  ) as favorited,
		  (count(distinct f.user_id)) as favorites_count,
		  u.username, u.bio, u.image,
		  (select exists
			(select 1 from follows where a.author_id = following_id and follower_id::text = $1)
		  ) as following
		from articles a
		left join users u on a.author_id = u.id
		left join article_tags at on at.article_id = a.id
		left join tags t on at.tag_id = t.id
		left join favorites f on f.article_id = a.id
		where a.deleted_at is null and slug = $2
		group by a.id, u.id;
	`

	/*
		example output:
		 slug | title |         description         |         body         |          created_at          |          updated_at          |     tags     | favorited | favorites_count |    username    |      bio      |                     image                      | following
		------+-------+-----------------------------+----------------------+------------------------------+------------------------------+--------------+-----------+-----------------+----------------+---------------+------------------------------------------------+-----------
		 slug | slug  | description cannot be empty | body cannot be empty | 2025-10-05 05:23:47.80455+00 | 2025-10-05 05:23:47.80455+00 | {sao,tai,vi} | t         |               3 | minhhoccode111 | i like golang | https://www.w3schools.com/howto/img_avatar.png | t
	*/

	a := entity.ArticleDetail{}
	row := r.Pool.QueryRow(ctx, query, userId, slug)
	err := row.Scan(
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

func (r *ArticleRepo) GetList(
	ctx context.Context,
	isFeed bool,
	userId, tag, author, favorited string,
	limit, offset uint64,
) (articles []entity.ArticlePreview, total uint64, err error) {
	var query string
	var args []any
	// NOTE: can't use squirrel because compldex queries should look complex :)
	if isFeed {
		query = `
		select a.slug, a.title, a.description, a.created_at, a.updated_at,
		  (select exists
			(select 1 from favorites where user_id::text = $1 and article_id = a.id)
		  ) as favorited,
		  u.username, u.bio, u.image,
		  (select exists
			(select 1 from follows where follower_id::text = $1 and following_id = u.id)
		  ) as following,
		  coalesce(array_agg(distinct t.name) filter (where t.name is not null), '{}') as tags,
		  count(distinct f.user_id) as favorites_count,
		  count(*) over() as articles_count
		from articles a
		left join users u on a.author_id = u.id
		left join article_tags at on at.article_id = a.id
		left join tags t on t.id = at.tag_id
		left join favorites f on f.article_id = a.id
		left join users u2 on f.user_id = u2.id
		where a.deleted_at is null
		  and (select exists
			(select 1 from follows where follower_id::text = $1
			  and following_id = u.id)
		  )
		group by a.id, u.id
		order by a.created_at desc
		limit $2
		offset $3;
	`
		args = []any{userId, limit, offset}
	} else {
		query = `
		select a.slug, a.title, a.description, a.created_at, a.updated_at,
		  (select exists
			(select 1 from favorites where user_id::text = $1 and article_id = a.id)
		  ) as favorited,
		  u.username, u.bio, u.image,
		  (select exists
			(select 1 from follows where follower_id::text = $1 and following_id = u.id)
		  ) as following,
		  coalesce(array_agg(distinct t.name) filter (where t.name is not null), '{}') as tags,
		  count(distinct f.user_id) as favorites_count,
		  count(*) over() as articles_count -- count all articles match before applying limit
		from articles a
		left join users u on a.author_id = u.id
		left join article_tags at on at.article_id = a.id
		left join tags t on t.id = at.tag_id
		left join favorites f on f.article_id = a.id
		left join users uf on f.user_id = uf.id
		where a.deleted_at is null
		  and ('' = $2 or u.username = $2) -- author, skip if empty
		  and ('' = $3 or uf.username = $3) -- favorited, skip if empty
		  and ('' = $4 or exists (select 1 from article_tags at2
			  left join tags t2 on at2.tag_id = t2.id
			  where at2.article_id = a.id and t2.name = $4)) -- tag, skip if empty
		group by a.id, u.id
		order by a.created_at desc
		limit $5
		offset $6;
	`
		args = []any{userId, author, favorited, tag, limit, offset}
	}

	/*
		example output:
		          slug           |         title         |         description         |          created_at           |          updated_at           | favorited | username | bio | image | following |     tags     | favorites_count | articles_count
		-------------------------+-----------------------+-----------------------------+-------------------------------+-------------------------------+-----------+----------+-----+-------+-----------+--------------+-----------------+----------------
		 title-cannot-be-empty-6 | title cannot be empty | description cannot be empty | 2025-10-02 13:38:00.168028+00 | 2025-10-02 13:38:00.168028+00 | t         | asd0     |     |       | t         | {tai,vi,sao} |               1 |              1
	*/

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ArticleRepo - GetList - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	articles = []entity.ArticlePreview{}
	for rows.Next() {
		var a entity.ArticlePreview
		err = rows.Scan(
			&a.Slug,
			&a.Title,
			&a.Description,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.Favorited,
			&a.Author.Username,
			&a.Author.Bio,
			&a.Author.Image,
			&a.Author.Following,
			&a.TagList,
			&a.FavoritesCount,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ArticleRepo - GetList - rows.Scan: %w", err)
		}
		articles = append(articles, a)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("ArticleRepo - GetList - rows.Err: %w", err)
	}

	return articles, total, nil
}
