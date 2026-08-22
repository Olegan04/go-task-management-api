package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager/internal/domain"
	"task-manager/internal/usecase"

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

func TestUserHandler_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	userUsecase := usecase.NewUserUsecase(mockRepo, "test-secret")

	handler := NewUserHandler(userUsecase)

	reqBody, _ := json.Marshal(domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.Register(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.UserResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "test@example.com", resp.Email)
	mockRepo.AssertExpectations(t)
}
