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

// Register -.
func (uc *UseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Register - util.HashPassword: %w",
			err,
		)
	}

	user.Password = hashedPassword
	user, err = uc.repo.StoreRegister(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Register - uc.repo.StoreRegister: %w",
			err,
		)
	}

	return user, nil
}

// Login -.
func (uc *UseCase) Login(ctx context.Context, loginCred entity.User) (entity.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, loginCred.Email)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Login - uc.repo.GetUserByEmail: %w",
			err,
		)
	}

	if !util.IsValidPassword(user.Password, loginCred.Password) {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Login - util.IsValidPassword: incorrect password",
		)
	}

	return user, nil
}

// Current -.
func (uc *UseCase) Current(ctx context.Context, userId string) (entity.User, error) {
	user, err := uc.repo.GetUserById(ctx, userId)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Current - uc.repo.GetUserById: %w",
			err,
		)
	}

	return user, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, user entity.User) (entity.User, error) {
	if user.Password != "" {
		hashedPassword, err := util.HashPassword(user.Password)
		if err != nil {
			return entity.User{}, fmt.Errorf(
				"UserUseCase - Register - util.HashPassword: %w",
				err,
			)
		}
		user.Password = hashedPassword
	}

	user, err := uc.repo.StoreRegister(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Register - uc.repo.StoreRegister: %w",
			err,
		)
	}

	return user, nil
}
