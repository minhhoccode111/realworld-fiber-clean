package tag

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase manages tag retrieval operations.
type UseCase struct {
	repo repo.TagRepo
}

// New constructs a tag use case with the provided repository.
func New(r repo.TagRepo) *UseCase {
	return &UseCase{repo: r}
}

// List returns tags with pagination across all articles.
func (uc *UseCase) List(ctx context.Context, limit, offset uint64,
) ([]entity.TagName, uint64, error) {
	tags, total, err := uc.repo.GetList(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"TagUseCase - GetTags - uc.repo.GetTags: %w",
			err,
		)
	}

	return tags, total, nil
}
