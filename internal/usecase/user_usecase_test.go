package usecase

import (
	"context"
	"task-manager/internal/domain"
	castomErr "task-manager/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserUsecase_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	usecase := NewUserUsecase(mockRepo, "test-secret")
	req := domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := usecase.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test@example.com", resp.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Register_EmailTaken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	existingUser := &domain.User{Email: "test@example.com"}
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	usecase := NewUserUsecase(mockRepo, "test-secret")
	req := domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := usecase.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, castomErr.ErrEmailTaken, err)
	mockRepo.AssertExpectations(t)
}
