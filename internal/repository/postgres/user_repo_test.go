package postgres

import (
	"context"
	"task-manager/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUserRepo_Create(t *testing.T) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)

	assert.NoError(t, err)
	t.Cleanup(func() {
		postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	dbConn, err := sqlx.Connect("postgres", dsn)
	assert.NoError(t, err)
	t.Cleanup(func() { dbConn.Close() })

	schema := `
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
    `
	_, err = dbConn.ExecContext(ctx, schema)
	assert.NoError(t, err)

	repo := NewUserRepo(dbConn)
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	err = repo.Create(ctx, user)
	assert.NoError(t, err)

	saved, err := repo.GetByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, saved.Email)
}
