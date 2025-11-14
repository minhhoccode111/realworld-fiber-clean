package favorite

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

type UseCase struct {
	repo repo.FavoriteRepo
}

func New(r repo.FavoriteRepo) *UseCase {
	return &UseCase{r}
}

func (uc *UseCase) Create(ctx context.Context, userID, slug string) error {
	err := uc.repo.StoreCreate(ctx, userID, slug)
	if err != nil {
		return fmt.Errorf("UseCase - Create - uc.repo.StoreCreate: %w", err)
	}

	return nil
}

func (uc *UseCase) Delete(ctx context.Context, userID, slug string) error {
	err := uc.repo.StoreDelete(ctx, userID, slug)
	if err != nil {
		return fmt.Errorf("UseCase - Delete - uc.repo.StoreDelete: %w", err)
	}

	return nil
}
