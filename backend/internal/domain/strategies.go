package domain

import "math"

// AddOperation implements Operation for addition.
// The empty struct{} means this type carries no data — it exists
// purely to attach the Execute method and satisfy the Operation interface.
type AddOperation struct{}

func (AddOperation) Execute(a, b float64) (float64, error) {
	return a + b, nil
}

// SubtractOperation implements Operation for subtraction.
type SubtractOperation struct{}

func (SubtractOperation) Execute(a, b float64) (float64, error) {
	return a - b, nil
}

// MultiplyOperation implements Operation for multiplication.
type MultiplyOperation struct{}

func (MultiplyOperation) Execute(a, b float64) (float64, error) {
	return a * b, nil
}

// DivideOperation implements Operation for division.
type DivideOperation struct{}

func (DivideOperation) Execute(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// PowerOperation implements Operation for exponentiation (a^b).
type PowerOperation struct{}

func (PowerOperation) Execute(a, b float64) (float64, error) {
	if a == 0 && b < 0 {
		return 0, ErrZeroToNegativePower
	}
	return math.Pow(a, b), nil
}

// SquareRootOperation implements Operation for square root.
// It's a unary operation, so b is intentionally unused — but it must
// still appear in the signature to satisfy the Operation interface.
// The underscore (_) tells both Go and any reader "yes, I know this
// parameter exists, and yes, I'm deliberately ignoring it."
type SquareRootOperation struct{}

func (SquareRootOperation) Execute(a, _ float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(a), nil
}

// PercentageOperation implements Operation for "a percent of b",
// e.g. Execute(20, 50) -> 10 (20% of 50 is 10).
type PercentageOperation struct{}

func (PercentageOperation) Execute(a, b float64) (float64, error) {
	return (a / 100) * b, nil
}