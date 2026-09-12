// This file exercises the Postgres adapter against a real database.
// Unlike the rest of the suite (domain, application, http), it cannot
// run without Postgres actually available, so it's gated behind the
// TEST_DATABASE_URL environment variable: unset (the default for
// `go test ./...` on a laptop with no database running), every test
// here calls t.Skip() and the rest of the suite is unaffected. Set it
// to a real connection string (see the README) or let CI provide one
// via a Postgres service container to exercise this file too.
package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/lib/pq"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/postgres"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// openTestDB connects to TEST_DATABASE_URL, (re-)applies the schema
// migration so the test starts from a known state regardless of what
// a previous run left behind, and registers cleanup for both the
// table contents and the connection itself.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping Postgres integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Drop first, then (re)create from the migration file: the
	// migration itself is CREATE TABLE IF NOT EXISTS, which is a no-op
	// against a table a previous test run already created — including
	// one left behind with an older schema. Dropping first guarantees
	// this run always starts from exactly what's in the file today.
	if _, err := db.Exec("DROP TABLE IF EXISTS calculations"); err != nil {
		t.Fatalf("failed to drop calculations table: %v", err)
	}
	applyMigration(t, db)

	return db
}

// applyMigration runs the actual migration file the project ships
// (migrations/001_create_calculations.sql) rather than a hand-copied
// CREATE TABLE statement, so this test breaks — loudly — the moment
// the real schema and the test's assumptions about it drift apart.
func applyMigration(t *testing.T, db *sql.DB) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate migration file relative to this test")
	}
	// this file:    backend/internal/adapters/persistence/postgres/repository_integration_test.go
	// migration:    backend/migrations/001_create_calculations.sql
	migrationPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "migrations", "001_create_calculations.sql")

	schema, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("failed to read migration file %q: %v", migrationPath, err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("failed to apply migration: %v", err)
	}
}

func TestRepository_SaveAndFindHistory(t *testing.T) {
	db := openTestDB(t)
	repo := postgres.NewRepository(db)
	ctx := context.Background()

	unary := domain.NewCalculation("11111111-1111-1111-1111-111111111111", domain.SquareRoot, 9, nil, 3)
	b := 4.0
	binary := domain.NewCalculation("22222222-2222-2222-2222-222222222222", domain.Addition, 10, &b, 14)

	if err := repo.Save(ctx, unary); err != nil {
		t.Fatalf("Save(unary) returned unexpected error: %v", err)
	}
	if err := repo.Save(ctx, binary); err != nil {
		t.Fatalf("Save(binary) returned unexpected error: %v", err)
	}

	history, err := repo.FindHistory(ctx, 10)
	if err != nil {
		t.Fatalf("FindHistory returned unexpected error: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("history length = %d, want 2", len(history))
	}

	// FindHistory orders newest first; `binary` was saved after `unary`.
	if history[0].ID != binary.ID {
		t.Errorf("history[0].ID = %v, want %v (newest first)", history[0].ID, binary.ID)
	}
	if history[0].OperandB == nil || *history[0].OperandB != b {
		t.Errorf("history[0].OperandB = %v, want %v", history[0].OperandB, b)
	}

	// The unary operation's OperandB must round-trip as NULL, not 0 —
	// that distinction is the entire reason OperandB is a *float64 in
	// the domain model to begin with.
	if history[1].OperandB != nil {
		t.Errorf("history[1].OperandB = %v, want nil (unary operation)", *history[1].OperandB)
	}
}

func TestRepository_FindHistory_RespectsLimit(t *testing.T) {
	db := openTestDB(t)
	repo := postgres.NewRepository(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		b := float64(i)
		calc := domain.NewCalculation(uuidFor(i), domain.Addition, 1, &b, 1+b)
		if err := repo.Save(ctx, calc); err != nil {
			t.Fatalf("Save returned unexpected error: %v", err)
		}
	}

	history, err := repo.FindHistory(ctx, 3)
	if err != nil {
		t.Fatalf("FindHistory returned unexpected error: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("history length = %d, want 3", len(history))
	}
}

// uuidFor generates a deterministic, distinct fake UUID per index —
// good enough for a primary key in a throwaway test table, without
// pulling in the uuid package just for test fixtures.
func uuidFor(i int) string {
	return "00000000-0000-0000-0000-00000000000" + string(rune('0'+i))
}
