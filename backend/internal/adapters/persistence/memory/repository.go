// Package memory provides an in-memory implementation of the
// application.CalculationRepository port. It exists purely to make
// the application layer testable in isolation, without a real
// database — this is one of the main payoffs of the hexagonal
// architecture's dependency inversion.
package memory

import (
	"context"
	"sync"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// This line does nothing at runtime — it's a compile-time assertion.
// It declares an unused variable (_) of type application.CalculationRepository,
// and tries to assign a nil *Repository to it. If Repository is
// missing any method the interface requires, or has a wrong method
// signature, this single line fails to compile with a clear message
// pointing exactly here — instead of a confusing error deep inside
// main.go, wherever the mismatch actually gets used.
var _ application.CalculationRepository = (*Repository)(nil)

// Repository is a thread-safe, in-memory CalculationRepository.
// It implicitly satisfies application.CalculationRepository — Go has
// no "implements" keyword; a type satisfies an interface automatically
// as soon as it has the right methods.
type Repository struct {
	mu    sync.Mutex
	items []domain.Calculation
}

// NewRepository creates an empty in-memory repository.
func NewRepository() *Repository {
	return &Repository{items: []domain.Calculation{}}
}

// Save appends a calculation to the in-memory slice.
func (r *Repository) Save(ctx context.Context, calc domain.Calculation) error {
	// sync.Mutex protects `items` from concurrent access. Go's HTTP
	// server handles each request in its own goroutine (lightweight
	// thread), so two requests could call Save() at the "same time" —
	// without locking, that's a data race that can corrupt the slice.
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, calc)
	return nil
}

// FindHistory returns up to `limit` calculations, most recent first.
func (r *Repository) FindHistory(ctx context.Context, limit int) ([]domain.Calculation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Save() only ever appends, so r.items is already in chronological
	// (oldest-first) order by construction — reversing it is enough to
	// get "most recent first", and unlike sorting by CreatedAt it can
	// never produce a wrong order from two calculations landing in the
	// same clock tick (time.Now() has limited resolution on some
	// platforms, and sort.Slice isn't even guaranteed stable on ties).
	n := len(r.items)
	reversed := make([]domain.Calculation, n)
	for i, calc := range r.items {
		reversed[n-1-i] = calc
	}

	if limit > 0 && limit < len(reversed) {
		reversed = reversed[:limit]
	}
	return reversed, nil
}
