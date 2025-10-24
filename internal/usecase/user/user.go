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
func (uc *UseCase) Register(ctx context.Context, userDTO entity.User) (entity.User, error) {
	hashedPassword, err := util.HashPassword(userDTO.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Register - util.HashPassword: %w",
			err,
		)
	}

	userDTO.Password = hashedPassword
	user, err := uc.repo.StoreRegister(ctx, userDTO)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Register - uc.repo.StoreRegister: %w",
			err,
		)
	}

	return user, nil
}

// Login -.
func (uc *UseCase) Login(ctx context.Context, userDTO entity.User) (entity.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"UserUseCase - Login - uc.repo.GetUserByEmail: %w",
			err,
		)
	}

	if !util.IsValidPassword(user.Password, userDTO.Password) {
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
		return entity.User{}, fmt.Errorf("UserUseCase - Current - uc.repo.GetUserById: %w", err)
	}

	return user, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, userDTO entity.User) (entity.User, error) {
	user, err := uc.repo.GetUserById(ctx, userDTO.Id)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Update - uc.repo.GetUserById: %w", err)
	}

	if userDTO.Password != "" {
		hashedPassword, err := util.HashPassword(userDTO.Password)
		if err != nil {
			return entity.User{}, fmt.Errorf("UserUseCase - Update - util.HashPassword: %w", err)
		}
		user.Password = hashedPassword
	}

	if userDTO.Email != "" {
		user.Email = userDTO.Email
	}

	if userDTO.Username != "" {
		user.Username = userDTO.Username
	}

	if userDTO.Bio != "" {
		user.Bio = userDTO.Bio
	}

	if userDTO.Image != "" {
		user.Image = userDTO.Image
	}

	err = uc.repo.StoreUpdate(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Update - uc.repo.StoreUpdate: %w", err)
	}

	return user, nil
}
