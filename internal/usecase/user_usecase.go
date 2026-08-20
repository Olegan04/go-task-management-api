package usecase

import (
	"context"
	"fmt"
	"task-manager/internal/domain"
	castomErr "task-manager/pkg/errors"
	"task-manager/pkg/utils"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserUsecase struct {
	repo      UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo UserRepository, secret string) *UserUsecase {
	return &UserUsecase{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (u *UserUsecase) Register(ctx context.Context, req domain.RegisterRequest) (*domain.UserResponse, error) {
	existing, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("usecase.Register: %w", err)
	}
	if existing != nil {
		return nil, castomErr.ErrEmailTaken
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("usecase.Register: hash failed: %w", err)
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

func (u *UserUsecase) Login(ctx context.Context, email, password string) (string, error) {
	existing, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("usecase.Login: %w", err)
	}
	if existing == nil {
		return "", castomErr.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(password)); err != nil {
		return "", castomErr.ErrInvalidCredentials
	}
	token, err := utils.GenerateToken(existing.ID, u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("usecase.Login: %w", err)
	}
	return token, nil
}
