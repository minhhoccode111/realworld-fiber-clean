package profile

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.ProfileRepo
}

// New -.
func New(r repo.ProfileRepo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) Detail(
	ctx context.Context,
	userId, username string,
) (entity.ProfilePreview, error) {
	profile, err := uc.repo.GetDetail(ctx, userId, username)
	if err != nil {
		return entity.ProfilePreview{}, fmt.Errorf(
			"ProfileUseCase - Detail - uc.repo.GetDetail: %w",
			err,
		)
	}

	return profile, nil
}

func (uc *UseCase) Follow(ctx context.Context, userId, username string) error {
	err := uc.repo.StoreCreate(ctx, userId, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Follow - uc.repo.StoreCreate: %w", err)
	}

	return nil
}
func (uc *UseCase) Unfollow(ctx context.Context, userId, username string) error {
	err := uc.repo.StoreDelete(ctx, userId, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Unfollow - uc.repo.StoreDelete: %w", err)
	}

	return nil
}
