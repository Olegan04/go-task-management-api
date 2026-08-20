package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	deliveryHttp "task-manager/internal/delivery/http"
	"task-manager/internal/repository/postgres"
	"task-manager/internal/usecase"
	"task-manager/pkg/db"
	"task-manager/pkg/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var jwtSecret = "Sheregesh=sci&&bike"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization is missing", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		userID, err := utils.ValidateToken(tokenString, []byte(jwtSecret))
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	dbCfg := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5433"),
		User:     getEnv("DB_USER", "app_user"),
		Password: getEnv("DB_PASSWORD", "app_user"),
		DBName:   getEnv("DB_Name", "task_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	database, err := db.NewPostgresDB(dbCfg)
	if err != nil {
		log.Fatalf("Fatal: cannot connect to db: %v", err)
	}
	defer database.Close()

	userRepo := postgres.NewUserRepo(database)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	userHandler := deliveryHttp.NewUserHandler(userUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Post("/api/register", userHandler.Register)
	r.Post("/api/login", userHandler.Login)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server is starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
