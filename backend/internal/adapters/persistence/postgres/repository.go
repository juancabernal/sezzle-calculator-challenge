// Package postgres provides the production implementation of the
// application.CalculationRepository port, backed by a real
// PostgreSQL database via the standard database/sql package.
package postgres

import (
	"context"
	"database/sql"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
	_ "github.com/lib/pq" // registers the "postgres" driver with database/sql
)

var _ application.CalculationRepository = (*Repository)(nil)

type Repository struct {
	db *sql.DB
}

// NewRepository wraps an already-open *sql.DB connection. The caller
// (main.go) owns opening and eventually closing that connection —
// this package only knows how to run queries against it.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, calc domain.Calculation) error {
	const query = `
		INSERT INTO calculations (id, operation, operand_a, operand_b, result, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		calc.ID, string(calc.Operation), calc.OperandA, calc.OperandB, calc.Result, calc.CreatedAt,
	)
	return err
}

func (r *Repository) FindHistory(ctx context.Context, limit int) ([]domain.Calculation, error) {
	const query = `
		SELECT id, operation, operand_a, operand_b, result, created_at
		FROM calculations
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Calculation
	for rows.Next() {
		var calc domain.Calculation
		var operation string

		if err := rows.Scan(
			&calc.ID, &operation, &calc.OperandA, &calc.OperandB, &calc.Result, &calc.CreatedAt,
		); err != nil {
			return nil, err
		}
		calc.Operation = domain.OperationType(operation)
		results = append(results, calc)
	}

	// rows.Err() catches errors that happened DURING iteration (e.g.
	// the connection dropped mid-scan) — distinct from the error
	// QueryContext itself might have returned before iteration started.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
