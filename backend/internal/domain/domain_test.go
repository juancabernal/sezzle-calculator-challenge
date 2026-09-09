package domain

import (
	"errors"
	"math"
	"testing"
)

// almostEqual compares floats with a small tolerance. Comparing
// float64 with == is risky due to rounding errors (e.g. math.Sqrt
// or math.Pow can produce 2.9999999999 instead of exactly 3).
func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

// TestOperations is a table-driven test: each entry describes one
// scenario (which operation, which inputs, what we expect back).
// The loop at the bottom runs every entry through the exact same logic.
func TestOperations(t *testing.T) {
	tests := []struct {
		name       string // human-readable subtest name
		opType     OperationType
		a, b       float64
		wantResult float64
		wantErr    error // nil if we expect success
	}{
		{name: "addition", opType: Addition, a: 2, b: 3, wantResult: 5},
		{name: "subtraction", opType: Subtraction, a: 10, b: 4, wantResult: 6},
		{name: "multiplication", opType: Multiplication, a: 3, b: 4, wantResult: 12},
		{name: "division", opType: Division, a: 10, b: 2, wantResult: 5},
		{name: "division by zero", opType: Division, a: 10, b: 0, wantErr: ErrDivisionByZero},
		{name: "power", opType: Power, a: 2, b: 3, wantResult: 8},
		{name: "zero to negative power", opType: Power, a: 0, b: -1, wantErr: ErrZeroToNegativePower},
		{name: "square root", opType: SquareRoot, a: 9, wantResult: 3},
		{name: "square root of negative", opType: SquareRoot, a: -4, wantErr: ErrNegativeSqrt},
		{name: "percentage", opType: Percentage, a: 20, b: 50, wantResult: 10},
	}

	for _, tc := range tests {
		// t.Run creates a "subtest": if "division by zero" fails,
		// `go test` tells you exactly that one failed, not just
		// "TestOperations failed" with no clue which case.
		t.Run(tc.name, func(t *testing.T) {
			op, err := NewOperation(tc.opType)
			if err != nil {
				t.Fatalf("NewOperation(%q) returned unexpected error: %v", tc.opType, err)
			}

			result, err := op.Execute(tc.a, tc.b)

			if tc.wantErr != nil {
				// errors.Is checks against our sentinel errors, not
				// against error message strings — this keeps tests
				// resilient even if we reword an error message later.
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("Execute(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute(%v, %v) returned unexpected error: %v", tc.a, tc.b, err)
			}
			if !almostEqual(result, tc.wantResult) {
				t.Errorf("Execute(%v, %v) = %v, want %v", tc.a, tc.b, result, tc.wantResult)
			}
		})
	}
}

// TestNewOperation_UnknownType verifies the factory rejects invalid
// operation types cleanly, instead of panicking or silently guessing.
func TestNewOperation_UnknownType(t *testing.T) {
	_, err := NewOperation(OperationType("not-a-real-operation"))
	if !errors.Is(err, ErrUnknownOperation) {
		t.Errorf("NewOperation with invalid type: error = %v, want %v", err, ErrUnknownOperation)
	}
}
