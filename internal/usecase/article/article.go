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

	err := uc.repo.StoreCreate(ctx, dto, tags)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Create - uc.repo.StoreCreate: %w",
			err,
		)
	}

	article, err := uc.repo.GetDetailBySlug(ctx, dto.AuthorId, dto.Slug)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Create - uc.repo.GetDetailBySlug: %w",
			err,
		)
	}

	return article, nil
}

func (uc *UseCase) List(ctx context.Context, isFeed bool, userId, tag, author, favorited string,
	limit, offset uint64) ([]entity.ArticlePreview, uint64, error) {
	articles, total, err := uc.repo.GetList(
		ctx,
		isFeed,
		userId,
		tag,
		author,
		favorited,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ArticleUseCase - List - uc.repo.GetList: %w",
			err,
		)
	}

	return articles, total, nil
}

func (uc *UseCase) Detail(ctx context.Context, userId, slug string) (entity.ArticleDetail, error) {
	article, err := uc.repo.GetDetailBySlug(ctx, userId, slug)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Detail - uc.repo.GetDetailBySlug: %w",
			err,
		)
	}

	return article, nil
}
