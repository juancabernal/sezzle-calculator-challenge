package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// CalculatorService is the main use case of the application: given an
// operation and its operands, compute the result via the domain layer
// and persist the outcome via the CalculationRepository port.
//
// Notice this struct holds an interface (CalculationRepository), not
// a concrete type like a Postgres client. That's constructor injection:
// whoever creates a CalculatorService decides which concrete adapter
// to plug in — the service itself never knows or cares.
type CalculatorService struct {
	repo CalculationRepository
}

// NewCalculatorService is the constructor. This explicit constructor
// pattern (instead of letting callers build the struct literal
// directly) is idiomatic Go for anything with dependencies — it makes
// "what does this type need to work?" obvious at a single call site.
func NewCalculatorService(repo CalculationRepository) *CalculatorService {
	return &CalculatorService{repo: repo}
}

// Calculate runs one arithmetic operation and stores the result.
// operandB is a pointer because unary operations (SquareRoot) don't
// have a second operand — see domain.Calculation for why that matters.
func (s *CalculatorService) Calculate(ctx context.Context, opType domain.OperationType, operandA float64, operandB *float64) (domain.Calculation, error) {
	operation, err := domain.NewOperation(opType)
	if err != nil {
		return domain.Calculation{}, err
	}

	// If operandB is nil (unary operation), pass 0 to Execute — the
	// concrete strategy (e.g. SquareRootOperation) already ignores
	// its second parameter, so this is safe.
	var b float64
	if operandB != nil {
		b = *operandB
	}

	result, err := operation.Execute(operandA, b)
	if err != nil {
		return domain.Calculation{}, err
	}

	calc := domain.NewCalculation(uuid.NewString(), opType, operandA, operandB, result)

	if err := s.repo.Save(ctx, calc); err != nil {
		return domain.Calculation{}, err
	}

	return calc, nil
}

// History returns the most recent calculations.
func (s *CalculatorService) History(ctx context.Context, limit int) ([]domain.Calculation, error) {
	return s.repo.FindHistory(ctx, limit)
}
