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
	userID, username string,
) (*entity.ProfilePreview, error) {
	p, err := uc.repo.GetDetail(ctx, userID, username)
	if err != nil {
		return nil, fmt.Errorf(
			"ProfileUseCase - Detail - uc.repo.GetDetail: %w",
			err,
		)
	}

	return p, nil
}

func (uc *UseCase) Follow(ctx context.Context, userID, username string) error {
	// NOTE: don't add concurrency because it's random between NoRows and NoEffect
	err := uc.repo.IsExisted(ctx, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Follow - uc.repo.IsExisted: %w", err)
	}

	err = uc.repo.StoreCreate(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Follow - uc.repo.StoreCreate: %w", err)
	}

	return nil
}

func (uc *UseCase) Unfollow(ctx context.Context, userID, username string) error {
	// NOTE: don't add concurrency because it's random between NoRows and NoEffect
	err := uc.repo.IsExisted(ctx, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Follow - uc.repo.IsExisted: %w", err)
	}

	err = uc.repo.StoreDelete(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("ProfileUseCase - Unfollow - uc.repo.StoreDelete: %w", err)
	}

	return nil
}
