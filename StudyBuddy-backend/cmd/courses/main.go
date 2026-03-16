package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/db"
	"studybuddy/backend/services/courses/delivery"
	"studybuddy/backend/services/courses/repository"
	"studybuddy/backend/services/courses/usecase"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	port := getEnv("COURSES_SERVER_PORT", "8082")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	dsn := getEnv("DATABASE_URL", "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPgCourseRepository(pool)
	svc := usecase.NewService(repo)

	handler := &delivery.CoursesHandler{
		List:   svc,
		Get:    svc,
		Create: svc,
		Update: svc,
		Delete: svc,
	}
	router := delivery.NewRouter(handler)

	// Protect mutating endpoints with JWT middleware.
	protect := auth.Middleware([]byte(jwtSecret))
	mux := http.NewServeMux()
	mux.Handle("/health", router)
	mux.Handle("/api/v1/courses", protect(http.HandlerFunc(handler.HandleCoursesCollection)))
	mux.Handle("/api/v1/courses/", protect(http.HandlerFunc(handler.HandleCourseItem)))

	log.Printf("courses service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
