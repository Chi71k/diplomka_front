package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	authjwt "studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/db"
	"studybuddy/backend/services/auth/delivery"
	"studybuddy/backend/services/auth/repository"
	"studybuddy/backend/services/auth/usecase"
	"time"
)

func main() {
	_ = godotenv.Load(".env")
	port := getEnv("AUTH_SERVER_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	jwtIssuer := getEnv("JWT_ISSUER", "studybuddy-auth")
	accessTTL := getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute)
	refreshTTL := getEnvDuration("JWT_REFRESH_TTL", 168*time.Hour)
	dsn := getEnv("DATABASE_URL", "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable")

	if len(jwtSecret) < 32 {
		log.Print("warning: JWT_SECRET should be at least 32 bytes for production")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewPgUserRepository(pool)
	hasher := usecase.PasswordAdapter{}
	jwtAdapter := usecase.JWTAdapter{
		Config: authjwt.Config{
			Secret:     []byte(jwtSecret),
			Issuer:     jwtIssuer,
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		},
	}

	registerUC := usecase.NewRegister(userRepo, hasher, jwtAdapter)
	loginUC := usecase.NewLogin(userRepo, hasher, jwtAdapter)

	handler := &delivery.AuthHandler{
		Register: registerUC,
		Login:    loginUC,
	}
	router := delivery.NewRouter(handler)

	log.Printf("auth service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}
