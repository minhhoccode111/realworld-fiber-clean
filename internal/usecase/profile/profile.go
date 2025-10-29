package profile

import (
	"context"

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

func (uc *UseCase) Detail(ctx context.Context)   {}
func (uc *UseCase) Follow(ctx context.Context)   {}
func (uc *UseCase) Unfollow(ctx context.Context) {}
