package article

import (
	"context"
	"errors"
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

func (uc *UseCase) Update(
	ctx context.Context,
	userId, oldSlug string,
	dto entity.Article,
) (entity.ArticleDetail, error) {
	a, err := uc.repo.GetBasicBySlug(ctx, oldSlug)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Update - uc.repo.GetBasicBySlug: %w",
			err,
		)
	}

	if a.AuthorId != userId {
		return entity.ArticleDetail{}, errors.New(
			"ArticleUseCase - Update - uc.repo.GetBasicBySlug: forbidden",
		)
	}

	baseSlug := slug.Make(dto.Title)
	a.Slug = baseSlug
	for i := 0; ; i++ {
		yes, err := uc.repo.CanSlugBeUsed(ctx, a.Id, a.Slug)
		if err != nil {
			return entity.ArticleDetail{}, fmt.Errorf(
				"ArticleUseCase - Update - uc.repo.CanSlugBeUsed: %w",
				err,
			)
		}

		if yes {
			break
		}

		a.Slug = baseSlug + "-" + strconv.Itoa(i)
	}

	if dto.Title != "" {
		a.Title = dto.Title
	}

	if dto.Body != "" {
		a.Body = dto.Body
	}

	if dto.Description != "" {
		a.Description = dto.Description
	}

	err = uc.repo.StoreUpdate(ctx, a)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Update - uc.repo.StoreUpdate: %w",
			err,
		)
	}

	article, err := uc.repo.GetDetailBySlug(ctx, userId, a.Slug)
	if err != nil {
		return entity.ArticleDetail{}, fmt.Errorf(
			"ArticleUseCase - Update - uc.repo.GetDetailBySlug: %w",
			err,
		)
	}

	return article, nil
}

func (uc *UseCase) Delete(ctx context.Context, userId, slug string) error {
	err := uc.repo.StoreDelete(ctx, userId, slug)
	if err != nil {
		return fmt.Errorf(
			"ArticleUseCase - Delete - uc.repo.StoreDelete: %w",
			err,
		)
	}

	return nil
}
