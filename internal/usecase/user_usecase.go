package usecase

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/domain"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(ctx context.Context, req domain.RegisterRequest) (*domain.UserResponse, error) {
	existing, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("usecase.Register: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email already taken")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("usecase.Register: hash falid: %w", err)
	}
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashed),
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("usecase.Register: %w", err)
	}
	return &domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
