package persistent

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

type ProfileRepo struct {
	*postgres.Postgres
}

func NewProfileRepo(pg *postgres.Postgres) *ProfileRepo {
	return &ProfileRepo{pg}
}

func (r *ProfileRepo) GetDetail(
	ctx context.Context,
	userId, username string,
) (entity.ProfilePreview, error) {
	panic("")
}

func (r *ProfileRepo) StoreCreate(ctx context.Context) {}

func (r *ProfileRepo) StoreDelete(ctx context.Context) {}
