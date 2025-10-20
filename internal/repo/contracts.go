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
		GetHistoryClone(
			ctx context.Context, limit, offset uint64,
		) (translations []entity.TranslationClone, translationsCount uint64, err error)
	}

	// TranslationCloneWebAPI -.
	TranslationCloneWebAPI interface {
		DoTranslate(entity.TranslationClone) (entity.TranslationClone, error)
	}
)
