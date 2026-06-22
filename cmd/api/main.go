package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/improver2108/jwt-login/db"
	"github.com/improver2108/jwt-login/db/sqlc"
	"github.com/improver2108/jwt-login/internal/cache"
	"github.com/improver2108/jwt-login/internal/handler"
	"github.com/improver2108/jwt-login/internal/repository"
	"github.com/improver2108/jwt-login/internal/routes"
	"github.com/improver2108/jwt-login/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()
	queries := sqlc.New(dbPool)

	redis := db.NewRedis()
	jtiCache := cache.NewJTICache(redis)

	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, jtiCache)

	r := chi.NewRouter()
	routes.RegisterUserRoutes(r, userHandler)
	routes.RegisterAuthRoutes(r, authHandler)
	fmt.Println("service running on :8082")
	http.ListenAndServe(":8082", r)
}
