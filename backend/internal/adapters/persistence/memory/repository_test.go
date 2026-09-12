package memory_test

import (
	"context"
	"testing"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/memory"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

func TestRepository_FindHistory_NewestFirst(t *testing.T) {
	repo := memory.NewRepository()
	ctx := context.Background()

	first := domain.NewCalculation("1", domain.Addition, 1, nil, 1)
	second := domain.NewCalculation("2", domain.Addition, 2, nil, 2)
	third := domain.NewCalculation("3", domain.Addition, 3, nil, 3)

	for _, calc := range []domain.Calculation{first, second, third} {
		if err := repo.Save(ctx, calc); err != nil {
			t.Fatalf("Save returned unexpected error: %v", err)
		}
	}

	history, err := repo.FindHistory(ctx, 10)
	if err != nil {
		t.Fatalf("FindHistory returned unexpected error: %v", err)
	}

	// Save() only ever appends; FindHistory must hand back the reverse
	// of insertion order ("most recent first") regardless of what
	// CreatedAt happens to hold for each entry.
	want := []string{third.ID, second.ID, first.ID}
	if len(history) != len(want) {
		t.Fatalf("history length = %d, want %d", len(history), len(want))
	}
	for i, id := range want {
		if history[i].ID != id {
			t.Errorf("history[%d].ID = %v, want %v", i, history[i].ID, id)
		}
	}
}

func TestRepository_FindHistory_RespectsLimit(t *testing.T) {
	repo := memory.NewRepository()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		calc := domain.NewCalculation(string(rune('a'+i)), domain.Addition, float64(i), nil, float64(i))
		if err := repo.Save(ctx, calc); err != nil {
			t.Fatalf("Save returned unexpected error: %v", err)
		}
	}

	history, err := repo.FindHistory(ctx, 2)
	if err != nil {
		t.Fatalf("FindHistory returned unexpected error: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("history length = %d, want 2", len(history))
	}
}

func TestRepository_FindHistory_EmptyRepository(t *testing.T) {
	repo := memory.NewRepository()

	history, err := repo.FindHistory(context.Background(), 10)
	if err != nil {
		t.Fatalf("FindHistory returned unexpected error: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("history length = %d, want 0", len(history))
	}
}
