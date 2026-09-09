// These types mirror EXACTLY the JSON shapes defined in the Go backend
// (backend/internal/adapters/http/dto.go). Keeping them in one file,
// named to match their backend counterpart, makes the contract between
// frontend and backend explicit and easy to keep in sync.

// Mirrors domain.OperationType from the backend. Using a union of
// string literals (instead of a plain `string`) means TypeScript will
// flag a typo like "ad" instead of "add" as a compile error, not a
// runtime bug discovered when the API rejects the request.
export type OperationType =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "power"
  | "sqrt"
  | "percentage";

// Mirrors CalculateRequest (Go). operandB is optional (`?`) because
// unary operations like sqrt don't send a second operand at all.
export interface CalculateRequest {
  operation: OperationType;
  operand_a: number;
  operand_b?: number;
}

// Mirrors CalculateResponse (Go).
export interface CalculateResponse {
  id: string;
  operation: OperationType;
  operand_a: number;
  operand_b?: number;
  result: number;
  created_at: string;
}

// Mirrors ErrorResponse (Go) — what the backend sends back on a 4xx/5xx.
export interface ApiErrorResponse {
  error: string;
}