// Package domain contains the core business rules of the calculator:
// what operations exist, what a "calculation" is, and the pure math
// itself. Nothing in this package knows about HTTP, JSON, or databases.
package domain

import "errors"

// OperationType identifies which arithmetic operation to perform.
// Using a distinct type (instead of a plain string) means the compiler
// stops us from passing an arbitrary, unchecked string wherever an
// OperationType is expected.
type OperationType string

const (
	Addition       OperationType = "add"
	Subtraction    OperationType = "subtract"
	Multiplication OperationType = "multiply"
	Division       OperationType = "divide"
	Power          OperationType = "power"
	SquareRoot     OperationType = "sqrt"
	Percentage     OperationType = "percentage"
)

// Domain-level sentinel errors. Callers (the application layer, the
// HTTP layer) can check against these with errors.Is() to decide how
// to respond — e.g. a division by zero should become an HTTP 400,
// not a generic 500.
var (
	ErrDivisionByZero      = errors.New("division by zero")
	ErrNegativeSqrt        = errors.New("cannot compute square root of a negative number")
	ErrZeroToNegativePower = errors.New("zero cannot be raised to a negative power")
	ErrUnknownOperation    = errors.New("unknown operation type")

	// ErrNonFiniteResult is returned when an otherwise-valid operation
	// produces a result that cannot be represented as JSON — NaN (e.g.
	// a negative base raised to a fractional power, like (-8)^0.5) or
	// +/-Inf (e.g. an exponent large enough to overflow float64). Both
	// are caught centrally in CalculatorService.Calculate rather than
	// in every individual strategy, since any operation involving very
	// large or otherwise unusual operands can in principle overflow.
	ErrNonFiniteResult = errors.New("this operation does not produce a finite real number for the given operands")
)

// Operation is the Strategy interface: every arithmetic operation
// (Addition, Division, SquareRoot, ...) implements this same method.
// Note that Execute always takes two operands — for unary operations
// like SquareRoot, the second operand (b) is simply ignored by that
// specific implementation. This keeps the interface uniform, which is
// what lets the factory (next step) treat every operation the same way.
type Operation interface {
	Execute(a, b float64) (float64, error)
}
