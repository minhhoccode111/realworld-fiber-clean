// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// TranslationCloneRepo -.
	TranslationCloneRepo interface {
		StoreTranslation(context.Context, entity.TranslationClone) error
		GetHistoryClone(ctx context.Context, limit, offset uint64,
		) (translations []entity.TranslationClone, total uint64, err error)
	}

	// TranslationCloneWebAPI -.
	TranslationCloneWebAPI interface {
		DoTranslate(entity.TranslationClone) (entity.TranslationClone, error)
	}

	// UserRepo -.
	UserRepo interface {
		StoreRegister(context.Context, entity.User) (entity.User, error)
		GetUserByEmail(context.Context, string) (entity.User, error)
		GetUserById(context.Context, string) (entity.User, error)
		StoreUpdate(context.Context, entity.User) error
	}

	// ArticleRepo -.
	ArticleRepo interface {
		StoreCreate(ctx context.Context, dto entity.Article, tags []string) (slug string, err error)
		CanSlugBeUsed(ctx context.Context, articleId string, slug string) (bool, error)
		// GetBySlug(context.Context, string) (entity.ArticleDetail, error)
	}

	// TagRepo -.
	TagRepo interface {
		GetList(ctx context.Context, limit, offset uint64,
		) (tags []entity.TagName, total uint64, err error)
		StoreList(context.Context, []string) ([]string, error)
	}
)
