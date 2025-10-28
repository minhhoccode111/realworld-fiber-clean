package comment

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.CommentRepo
}

// New -.
func New(r repo.CommentRepo) *UseCase {
	return &UseCase{repo: r}
}

// Create -.
func (uc *UseCase) Create(
	ctx context.Context,
	slug string,
	dto entity.Comment,
) (entity.CommentDetail, error) {
	id, err := uc.repo.StoreCreate(ctx, slug, dto)
	if err != nil {
		return entity.CommentDetail{}, fmt.Errorf(
			"CommentUseCase - Create - uc.repo.StoreCreate: %w",
			err,
		)
	}

	comment, err := uc.repo.GetDetailById(ctx, dto.AuthorId, id)
	if err != nil {
		return entity.CommentDetail{}, fmt.Errorf(
			"CommentUseCase - Create - uc.repo.GetDetailById: %w",
			err,
		)
	}

	return comment, nil
}
