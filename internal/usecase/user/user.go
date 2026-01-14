package user

import (
	"context"
	"fmt"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utilities"
)

// UseCase handles business logic for user accounts.
type UseCase struct {
	repo repo.UserRepo
}

// New constructs a user use case with the provided repository.
func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

// Register hashes the password and stores a new user.
func (uc *UseCase) Register(ctx context.Context, u *entity.User) error {
	hashedPassword, err := utilities.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf(
			"UserUseCase - Register - utils.HashPassword: %w",
			err,
		)
	}

	u.Password = hashedPassword

	err = uc.repo.StoreRegister(ctx, u)
	if err != nil {
		return fmt.Errorf(
			"UserUseCase - Register - uc.repo.StoreRegister: %w",
			err,
		)
	}

	return nil
}

// Login authenticates a user by verifying the supplied credentials.
func (uc *UseCase) Login(ctx context.Context, dto *entity.User) (*entity.User, error) {
	u, err := uc.repo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, fmt.Errorf(
			"UserUseCase - Login - uc.repo.GetUserByEmail: %w",
			err,
		)
	}

	if !utilities.IsValidPassword(u.Password, dto.Password) {
		return nil, fmt.Errorf(
			"UserUseCase - Login - utils.IsValidPassword: %w",
			entity.ErrInvalidCredentials,
		)
	}

	return u, nil
}

// Current returns user details by ID.
func (uc *UseCase) Current(ctx context.Context, userID string) (*entity.User, error) {
	u, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Current - uc.repo.GetUserByID: %w", err)
	}

	return u, nil
}

// Update applies provided changes to a user and persists them.
func (uc *UseCase) Update(ctx context.Context, dto *entity.User) (*entity.User, error) {
	u, err := uc.repo.GetUserByID(ctx, dto.ID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Update - uc.repo.GetUserByID: %w", err)
	}

	if dto.Password != "" {
		hashedPassword, err := utilities.HashPassword(dto.Password)
		if err != nil {
			return nil, fmt.Errorf("UserUseCase - Update - utils.HashPassword: %w", err)
		}

		u.Password = hashedPassword
	}

	if dto.Email != "" {
		u.Email = dto.Email
	}

	if dto.Username != "" {
		u.Username = dto.Username
	}

	if dto.Bio != "" {
		u.Bio = dto.Bio
	}

	if dto.Image != "" {
		u.Image = dto.Image
	}

	err = uc.repo.StoreUpdate(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Update - uc.repo.StoreUpdate: %w", err)
	}

	return u, nil
}
