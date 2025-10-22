package tag

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.TagRepo
}

// New -.
func New(r repo.TagRepo) *UseCase {
	return &UseCase{repo: r}
}

// GetTags - get all tags of all articles
func (uc *UseCase) GetTags(ctx context.Context, limit, offset uint64,
) ([]entity.Tag, uint64, error) {
	tags, total, err := uc.repo.RetrieveTags(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"TagUseCase - GetTags - uc.repo.GetTags: %w",
			err,
		)
	}

	return tags, total, nil
}
