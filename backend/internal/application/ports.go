// Package application contains the use cases (business workflows) of
// the calculator. It orchestrates the domain layer and talks to the
// outside world only through the ports (interfaces) defined here —
// it never imports a concrete database driver or HTTP framework.
package application

import (
	"context"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// CalculationRepository is a port: a contract for persisting and
// retrieving Calculations. The application layer depends on THIS
// interface, never on a concrete database implementation. Whoever
// wires the app together (main.go) decides which concrete adapter
// (Postgres, in-memory, ...) satisfies it.
type CalculationRepository interface {
	// Save persists a single calculation.
	Save(ctx context.Context, calc domain.Calculation) error

	// FindHistory returns the most recent calculations, newest first,
	// limited to at most `limit` results.
	FindHistory(ctx context.Context, limit int) ([]domain.Calculation, error)
}
