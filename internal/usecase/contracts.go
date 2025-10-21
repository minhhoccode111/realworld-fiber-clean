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
		GetHistory(ctx context.Context, limit, offset uint64,
		) (translations []entity.TranslationClone, total uint64, err error)
	}

	// Tag -.
	Tag interface {
		GetTags(ctx context.Context, limit, offset uint64,
		) (tags []entity.Tag, total uint64, err error)
	}
)
