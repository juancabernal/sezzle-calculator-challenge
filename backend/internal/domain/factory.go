package domain

// NewOperation is a factory function: given an OperationType, it
// returns the concrete Operation strategy that knows how to execute it.
//
// This is the ONLY place in the entire codebase that maps a string-like
// value to a concrete operation. If we add an 8th operation later
// (say, Modulo), we add one case here and one new file for the
// strategy — nothing else in the application needs to change. That's
// the Open/Closed Principle (the "O" in SOLID) in practice.
func NewOperation(opType OperationType) (Operation, error) {
	switch opType {
	case Addition:
		return AddOperation{}, nil
	case Subtraction:
		return SubtractOperation{}, nil
	case Multiplication:
		return MultiplyOperation{}, nil
	case Division:
		return DivideOperation{}, nil
	case Power:
		return PowerOperation{}, nil
	case SquareRoot:
		return SquareRootOperation{}, nil
	case Percentage:
		return PercentageOperation{}, nil
	default:
		return nil, ErrUnknownOperation
	}
}
