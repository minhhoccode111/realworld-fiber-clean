package translation

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase coordinates translation operations using storage and external API.
type UseCase struct {
	repo   repo.TranslationRepo
	webAPI repo.TranslationWebAPI
}

// New constructs a translation use case with the given store and web API.
func New(r repo.TranslationRepo, w repo.TranslationWebAPI) *UseCase {
	return &UseCase{
		repo:   r,
		webAPI: w,
	}
}

// History retrieves stored translation history.
func (uc *UseCase) History(ctx context.Context) (entity.TranslationHistory, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf(
			"TranslationUseCase - History - s.repo.GetHistory: %w",
			err,
		)
	}

	return entity.TranslationHistory{History: translations}, nil
}

// Translate calls the web API then stores and returns the translation result.
func (uc *UseCase) Translate(
	ctx context.Context,
	t entity.Translation,
) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.webAPI.Translate: %w",
			err,
		)
	}

	err = uc.repo.Store(ctx, translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf(
			"TranslationUseCase - Translate - s.repo.Store: %w",
			err,
		)
	}

	return translation, nil
}
