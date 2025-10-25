package article

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.ArticleRepo
}

// New -.
func New(r repo.ArticleRepo) *UseCase {
	return &UseCase{repo: r}
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, dto entity.Article, tags []string,
) (entity.ArticleDetail, error) {
	// 0/ find usable slug
	baseSlug := slug.Make(dto.Title)
	dto.Slug = baseSlug
	for i := 0; ; i++ {
		yes, err := uc.repo.CanSlugBeUsed(ctx, "", dto.Slug)
		if err != nil {
			return entity.ArticleDetail{}, fmt.Errorf(
				"ArticleUseCase - Create - uc.repo.CanSlugBeUsed: %w",
				err,
			)
		}

		if yes {
			break
		}

		dto.Slug = baseSlug + "-" + strconv.Itoa(i)
	}

	// 1/ call article store create with new slug, and slice tags
	err := uc.repo.StoreCreate(ctx, dto, tags)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Create - uc.repo.StoreCreate: %w",
			err,
		)
	}

	// 2/ use new slug to retrieve article detail
	article, err := uc.repo.GetDetailBySlug(ctx, dto.AuthorId, dto.Slug)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Create - uc.repo.GetDetailBySlug: %w",
			err,
		)
	}

	return article, nil
}
