import type { OperationType } from "../types/calculation";

// Operations that only need one operand — the second input in the
// form should be disabled/ignored for these.
const UNARY_OPERATIONS: OperationType[] = ["sqrt"];

export function isUnaryOperation(operation: OperationType): boolean {
  return UNARY_OPERATIONS.includes(operation);
}

// validateInputs is a pure function: given the raw form state, it
// returns an error message to display, or null if everything is
// valid enough to submit. It never touches the DOM or the network —
// that's what makes it trivial to unit test in isolation.
export function validateInputs(
  operation: OperationType,
  operandA: string,
  operandB: string
): string | null {
  if (operandA.trim() === "") {
    return "Please enter the first number.";
  }
  if (isNaN(Number(operandA))) {
    return "The first number is not valid.";
  }

  if (!isUnaryOperation(operation)) {
    if (operandB.trim() === "") {
      return "Please enter the second number.";
    }
    if (isNaN(Number(operandB))) {
      return "The second number is not valid.";
    }
    if (operation === "divide" && Number(operandB) === 0) {
      return "Cannot divide by zero.";
    }
  }

  return null;
}