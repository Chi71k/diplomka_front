package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"studybuddy/backend/pkg/db"
	"studybuddy/backend/services/users/delivery"
	"studybuddy/backend/services/users/repository"
	"studybuddy/backend/services/users/usecase"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	port := getEnv("USERS_SERVER_PORT", "8081")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	dsn := getEnv("DATABASE_URL", "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	profileRepo := repository.NewPgProfileRepository(pool)
	_ = repository.NewPgInterestRepository(pool) // placeholder for future interests endpoints

	getMeUC := usecase.NewGetMe(profileRepo)
	updateMeUC := usecase.NewUpdateMe(profileRepo)
	deleteMeUC := usecase.NewDeleteMe(profileRepo)

	searchURL := getEnv("SEARCH_SERVICE_URL", "http://localhost:8083")

	handler := &delivery.UsersHandler{
		GetMe:            getMeUC,
		UpdateMe:         updateMeUC,
		DeleteMe:         deleteMeUC,
		SearchServiceURL: searchURL,
	}
	router := delivery.NewRouter(handler, []byte(jwtSecret))

	log.Printf("users service listening on :%s", port)
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
