package persistent

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type ArticleRepo struct {
	*postgres.Postgres
}

func NewArticleRepo(pg *postgres.Postgres) *ArticleRepo {
	return &ArticleRepo{pg}
}

func (r *ArticleRepo) StoreCreate(ctx context.Context, dto entity.Article, slugs []string,
) (string, error) {
	return "", nil
}

func (r *ArticleRepo) CanSlugBeUsed(ctx context.Context, articleId, slug string) (bool, error) {
	return false, nil
}
