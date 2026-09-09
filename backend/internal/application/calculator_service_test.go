package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/memory"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

func TestCalculatorService_Calculate_Success(t *testing.T) {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	ctx := context.Background()

	b := 4.0
	calc, err := service.Calculate(ctx, domain.Addition, 10, &b)
	if err != nil {
		t.Fatalf("Calculate returned unexpected error: %v", err)
	}
	if calc.Result != 14 {
		t.Errorf("Result = %v, want 14", calc.Result)
	}
	if calc.ID == "" {
		t.Error("expected a generated ID, got empty string")
	}

	history, err := service.History(ctx, 10)
	if err != nil {
		t.Fatalf("History returned unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("History length = %d, want 1", len(history))
	}
	if history[0].ID != calc.ID {
		t.Errorf("History[0].ID = %v, want %v", history[0].ID, calc.ID)
	}
}

func TestCalculatorService_Calculate_UnaryOperation(t *testing.T) {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	ctx := context.Background()

	calc, err := service.Calculate(ctx, domain.SquareRoot, 9, nil)
	if err != nil {
		t.Fatalf("Calculate returned unexpected error: %v", err)
	}
	if calc.Result != 3 {
		t.Errorf("Result = %v, want 3", calc.Result)
	}
	if calc.OperandB != nil {
		t.Errorf("OperandB = %v, want nil", *calc.OperandB)
	}
}

func TestCalculatorService_Calculate_DomainError(t *testing.T) {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	ctx := context.Background()

	zero := 0.0
	_, err := service.Calculate(ctx, domain.Division, 10, &zero)
	if !errors.Is(err, domain.ErrDivisionByZero) {
		t.Errorf("error = %v, want %v", err, domain.ErrDivisionByZero)
	}

	history, _ := service.History(ctx, 10)
	if len(history) != 0 {
		t.Errorf("History length = %d, want 0 (failed calculation should not persist)", len(history))
	}
}

func TestCalculatorService_Calculate_UnknownOperation(t *testing.T) {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	ctx := context.Background()

	b := 1.0
	_, err := service.Calculate(ctx, domain.OperationType("bogus"), 1, &b)
	if !errors.Is(err, domain.ErrUnknownOperation) {
		t.Errorf("error = %v, want %v", err, domain.ErrUnknownOperation)
	}
}

func TestCalculatorService_History_RespectsLimit(t *testing.T) {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		b := float64(i)
		if _, err := service.Calculate(ctx, domain.Addition, 1, &b); err != nil {
			t.Fatalf("Calculate returned unexpected error: %v", err)
		}
	}

	history, err := service.History(ctx, 3)
	if err != nil {
		t.Fatalf("History returned unexpected error: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("History length = %d, want 3", len(history))
	}
}
