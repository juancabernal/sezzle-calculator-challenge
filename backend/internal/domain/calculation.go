package domain

import "time"

// Calculation is the aggregate we persist: a record of one arithmetic
// operation that was performed, including its result and when it
// happened. This is what turns the calculator from a stateless
// function into something with real, queryable history.
type Calculation struct {
	ID        string
	Operation OperationType
	OperandA  float64
	// OperandB is a pointer (*float64), not a plain float64, because
	// unary operations like SquareRoot genuinely have no second
	// operand — nil means "not applicable", which is different from
	// "the operand happens to be zero". A plain float64 can't express
	// that distinction; its zero value (0) would be indistinguishable
	// from a real zero operand.
	OperandB  *float64
	Result    float64
	CreatedAt time.Time
}

// NewCalculation builds a Calculation with a generated ID and the
// current timestamp, so callers never have to worry about setting
// those fields themselves.
func NewCalculation(id string, opType OperationType, a float64, b *float64, result float64) Calculation {
	return Calculation{
		ID:        id,
		Operation: opType,
		OperandA:  a,
		OperandB:  b,
		Result:    result,
		CreatedAt: time.Now().UTC(),
	}
}
