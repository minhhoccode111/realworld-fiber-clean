package translationclone

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo   repo.TranslationCloneRepo
	webAPI repo.TranslationCloneWebAPI
}

// New -.
func New(r repo.TranslationCloneRepo, w repo.TranslationCloneWebAPI) *UseCase {
	return &UseCase{
		repo:   r,
		webAPI: w,
	}
}

// GetHistory - getting translate history from store
func (uc *UseCase) GetHistory(
	ctx context.Context,
	limit, offset uint64,
) ([]entity.TranslationClone, uint64, error) {
	translations, total, err := uc.repo.GetHistoryClone(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"TranslationCloneUseCase - GetHistory - uc.repo.GetHistoryClone: %w",
			err,
		)
	}

	return translations, total, nil
}

func (uc *UseCase) PostTranslate(
	ctx context.Context,
	t entity.TranslationClone,
) (entity.TranslationClone, error) {
	translation, err := uc.webAPI.DoTranslate(t)
	if err != nil {
		return entity.TranslationClone{}, fmt.Errorf(
			"TranslationCloneUseCase - PostTranslate - uc.webAPI.DoTranslate: %w",
			err,
		)
	}

	err = uc.repo.StoreTranslation(ctx, translation)
	if err != nil {
		return entity.TranslationClone{}, fmt.Errorf(
			"TranslationCloneUseCase - PostTranslate - uc.repo.StoreTranslation: %w",
			err,
		)
	}

	return translation, nil
}
