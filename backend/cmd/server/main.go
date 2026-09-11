package main

import (
	"database/sql"
	"log"
	"net/http"

	httpadapter "github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/http"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/memory"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/postgres"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/config"
)

func main() {
	cfg := config.Load()

	// 1. Build the concrete persistence adapter. This is the ONLY
	//    place in the whole codebase that decides in-memory vs.
	//    Postgres — everything downstream just sees
	//    application.CalculationRepository.
	repo, closeDB := buildRepository(cfg)
	defer closeDB()

	service := application.NewCalculatorService(repo)
	handler := httpadapter.NewHandler(service)
	router := httpadapter.NewRouter(handler)

	log.Printf("calculator backend listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// buildRepository picks the repository implementation based on
// configuration. It also returns a cleanup function — closeDB() —
// so main() can defer closing the database connection regardless of
// which branch was taken (the in-memory case just returns a no-op).
func buildRepository(cfg config.Config) (application.CalculationRepository, func()) {
	if cfg.DatabaseURL == "" {
		log.Println("no DATABASE_URL set, using in-memory repository")
		return memory.NewRepository(), func() {}
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	// sql.Open only validates arguments — it doesn't actually connect.
	// Ping forces a real connection attempt now, so a misconfigured
	// database fails loudly at startup, not on the first request.
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to Postgres, using postgres repository")
	return postgres.NewRepository(db), func() { db.Close() }
}
