package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"studybuddy/backend/pkg/crypto"
	"studybuddy/backend/pkg/db"
	pkggcal "studybuddy/backend/pkg/gcal"
	"studybuddy/backend/services/availability/delivery"
	"studybuddy/backend/services/availability/repository"
	"studybuddy/backend/services/availability/usecase"
)

func main() {
	_ = godotenv.Load(".env")

	port := getEnv("AVAILABILITY_SERVER_PORT", "8083")

	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	if len(jwtSecret) < 32 {
		log.Print("warning: JWT_SECRET should be at least 32 bytes for production")
	}

	encKey := crypto.MustKeyFromEnv("ENCRYPTION_KEY")

	stateKey := crypto.MustKeyFromEnv("GCAL_STATE_KEY")

	gcalClientID := getEnv("GCAL_CLIENT_ID", "")
	gcalClientSecret := getEnv("GCAL_CLIENT_SECRET", "")
	gcalRedirectURL := getEnv(
		"GCAL_REDIRECT_URL",
		"http://localhost:8083/api/v1/availability/gcal/callback",
	)

	if gcalClientID == "" || gcalClientSecret == "" {
		log.Print("warning: GCAL_CLIENT_ID or GCAL_CLIENT_SECRET not set — google calendar integration will not work")
	}

	dsn := getEnv("DATABASE_URL", "postgres://studybuddy:studybuddy@localhost:5432/studybuddy?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	slotRepo := repository.NewPgSlotRepository(pool)
	gcalRepo := repository.NewPgGCalRepository(pool, encKey)

	gcalProvider := pkggcal.New(pkggcal.Config{
		ClientID:     gcalClientID,
		ClientSecret: gcalClientSecret,
		RedirectURL:  gcalRedirectURL,
	})

	// ── use cases ─────────────────────────────────────────────────────────────
	listSlotsUC := usecase.NewListSlots(slotRepo)
	createSlotUC := usecase.NewCreateSlot(slotRepo)
	deleteSlotUC := usecase.NewDeleteSlot(slotRepo)
	gcalConnectUC := usecase.NewGCalConnect(gcalProvider, gcalRepo, stateKey)
	gcalImportUC := usecase.NewGCalImport(gcalProvider, gcalRepo, slotRepo)
	gcalDisconnectUC := usecase.NewGCalDisconnect(gcalRepo, slotRepo)

	// ── handler + router ─────────────────────────────────────────────────────
	handler := &delivery.AvailabilityHandler{
		ListSlots:      listSlotsUC,
		CreateSlot:     createSlotUC,
		DeleteSlot:     deleteSlotUC,
		GCalConnect:    gcalConnectUC,
		GCalImport:     gcalImportUC,
		GCalDisconnect: gcalDisconnectUC,
	}
	router := delivery.NewRouter(handler, []byte(jwtSecret))

	log.Printf("availability service listening on :%s", port)
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
