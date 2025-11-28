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
	dto *entity.Comment,
) (*entity.CommentDetail, error) {
	id, err := uc.repo.StoreCreate(ctx, slug, dto)
	if err != nil {
		return nil, fmt.Errorf(
			"CommentUseCase - Create - uc.repo.StoreCreate: %w",
			err,
		)
	}

	c, err := uc.repo.GetDetailByID(ctx, dto.AuthorID, id)
	if err != nil {
		return nil, fmt.Errorf(
			"CommentUseCase - Create - uc.repo.GetDetailByID: %w",
			err,
		)
	}

	return c, nil
}

func (uc *UseCase) List(
	ctx context.Context,
	userID, slug string,
	limit, offset uint64,
) ([]entity.CommentDetail, uint64, error) {
	comments, total, err := uc.repo.GetList(ctx, userID, slug, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"CommentUseCase - List - uc.repo.GetList: %w",
			err,
		)
	}

	return comments, total, nil
}

func (uc *UseCase) Delete(
	ctx context.Context,
	userID, slug, commentID string,
	userRole entity.Role,
) error {
	c, err := uc.repo.GetBasicByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf(
			"CommentUseCase - Delete - uc.repo.GetBasicByID: %w",
			err,
		)
	}

	if userRole != entity.AdminRole && c.AuthorID != userID {
		return fmt.Errorf(
			"CommentUseCase - Delete - uc.repo.GetBasicByID: %w",
			entity.ErrForbidden,
		)
	}

	err = uc.repo.StoreDelete(ctx, userID, slug, commentID)
	if err != nil {
		return fmt.Errorf(
			"CommentUseCase - Delete - uc.repo.StoreDelete: %w",
			err,
		)
	}

	return nil
}
