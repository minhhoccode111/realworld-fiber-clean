package favorite

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

type FavoriteUseCase struct {
	repo repo.FavoriteRepo
}

func New(r repo.FavoriteRepo) *FavoriteUseCase {
	return &FavoriteUseCase{r}
}

func (uc *FavoriteUseCase) Create(ctx context.Context, userId, slug string) error {
	err := uc.repo.StoreCreate(ctx, userId, slug)
	if err != nil {
		return fmt.Errorf("FavoriteUseCase - Create - uc.repo.StoreCreate: %w", err)
	}

	return nil
}

func (uc *FavoriteUseCase) Delete(ctx context.Context, userId, slug string) error {
	err := uc.repo.StoreDelete(ctx, userId, slug)
	if err != nil {
		return fmt.Errorf("FavoriteUseCase - Delete - uc.repo.StoreDelete: %w", err)
	}

	return nil
}
