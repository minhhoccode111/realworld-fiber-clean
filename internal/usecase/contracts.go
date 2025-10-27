// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) (entity.TranslationHistory, error)
	}

	// TranslationClone -.
	// Try to extend the Translation example with pagination
	TranslationClone interface {
		PostTranslate(context.Context, entity.TranslationClone) (entity.TranslationClone, error)
		GetHistory(
			ctx context.Context,
			limit, offset uint64,
		) (translations []entity.TranslationClone, total uint64, err error)
	}

	// User -.
	User interface {
		Register(context.Context, entity.User) (entity.User, error)
		Login(context.Context, entity.User) (entity.User, error)
		Current(context.Context, string) (entity.User, error)
		Update(context.Context, entity.User) (entity.User, error)
	}

	// Article -.
	Article interface {
		Create(context.Context, entity.Article, []string) (entity.ArticleDetail, error)
		List(
			ctx context.Context,
			isFeed bool,
			userId, tag, author, favorited string,
			limit, offset uint64,
		) ([]entity.ArticlePreview, uint64, error)
		Detail(ctx context.Context, userId string, slug string) (entity.ArticleDetail, error)
		Update(
			ctx context.Context,
			userId string,
			oldSlug string,
			dto entity.Article,
		) (entity.ArticleDetail, error)
	}

	// Tag -.
	Tag interface {
		List(
			ctx context.Context,
			limit, offset uint64,
		) (tags []entity.TagName, total uint64, err error)
	}
)
