package errors

import "errors"

var (
	ErrEmailTaken         = errors.New("user_usecase: email already taken")
	ErrInvalidCredentials = errors.New("user_usecase: invalid credentials")
	ErrInvalidStatus      = errors.New("task_usecase: invalid task status")
	ErrEmptyTitle         = errors.New("task_usecase: task title cannot be empty")
	ErrInvalidToken       = errors.New("jwt.ValidateToken: invalid token")
	ErrInvalidClaimType   = errors.New("jwt.ValidateToken: invalid claim type")
	ErrWithUserID         = errors.New("jwt.ValidateToken: missing or invalid user_id field")
	ErrTaskNotFound       = errors.New("task not found")
)
