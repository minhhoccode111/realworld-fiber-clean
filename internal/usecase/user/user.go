package user

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/util"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
}

// New -.
func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

// RegisterUser -.
func (uc *UseCase) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - RegisterUser - util.HashPassword: %w",
			err,
		)
	}

	user.Password = hashedPassword
	user, err = uc.repo.StoreRegisterUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - RegisterUser - uc.repo.StoreRegisterUser: %w",
			err,
		)
	}

	return user, nil
}
