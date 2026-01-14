package favorite

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase manages favorite actions for articles.
type UseCase struct {
	repo repo.FavoriteRepo
}

// New constructs a favorite use case with the provided repository.
func New(r repo.FavoriteRepo) *UseCase {
	return &UseCase{r}
}

// Create marks an article as favorited for the given user.
func (uc *UseCase) Create(ctx context.Context, userID, slug string) error {
	err := uc.repo.StoreCreate(ctx, userID, slug)
	if err != nil {
		return fmt.Errorf("UseCase - Create - uc.repo.StoreCreate: %w", err)
	}

	return nil
}

// Delete removes a favorite marker for the given user and article.
func (uc *UseCase) Delete(ctx context.Context, userID, slug string) error {
	err := uc.repo.StoreDelete(ctx, userID, slug)
	if err != nil {
		return fmt.Errorf("UseCase - Delete - uc.repo.StoreDelete: %w", err)
	}

	return nil
}
