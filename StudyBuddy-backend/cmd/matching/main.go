package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"studybuddy/backend/pkg/db"
	"studybuddy/backend/services/matching/delivery"
	"studybuddy/backend/services/matching/repository"
	"studybuddy/backend/services/matching/usecase"
)

func main() {
	_ = godotenv.Load(".env")

	port := getEnv("MATCHING_SERVER_PORT", "8084")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	dsn := getEnv("DATABASE_URL", "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	matchRepo := repository.NewPgMatchRepository(pool)
	candidateStore := repository.NewPgCandidateStore(pool)
	profileClient := repository.NewPgProfileClient(pool)
	slotClient := repository.NewPgSlotClient(pool)
	courseClient := repository.NewPgCourseClient(pool)

	listCandidatesUC := usecase.NewListCandidates(matchRepo, profileClient, slotClient, courseClient, candidateStore)
	sendRequestUC := usecase.NewSendMatchRequest(matchRepo)
	respondUC := usecase.NewRespondToMatch(matchRepo)
	cancelUC := usecase.NewCancelMatch(matchRepo)
	listMatchesUC := usecase.NewListMatches(matchRepo)

	handler := &delivery.MatchingHandler{
		ListCandidates: listCandidatesUC,
		SendRequest:    sendRequestUC,
		Respond:        respondUC,
		Cancel:         cancelUC,
		ListMatches:    listMatchesUC,
	}
	router := delivery.NewRouter(handler, []byte(jwtSecret))

	log.Printf("matching service listening on :%s", port)
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
