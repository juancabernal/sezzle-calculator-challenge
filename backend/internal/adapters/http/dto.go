// Package http is the inbound adapter: it translates HTTP requests
// into calls to the application layer, and translates the results
// back into JSON responses. It knows about domain and application,
// but domain and application never know this package exists.
package http

import "github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"

// CalculateRequest is the exact shape of the JSON body a client sends
// to POST /calculate. The `json:"..."` tags control how Go's
// encoding/json package maps struct fields to/from JSON keys.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	OperandA  float64  `json:"operand_a"`
	OperandB  *float64 `json:"operand_b,omitempty"`
}

// CalculateResponse is the exact shape of the JSON we send back after
// a successful calculation.
type CalculateResponse struct {
	ID        string   `json:"id"`
	Operation string   `json:"operation"`
	OperandA  float64  `json:"operand_a"`
	OperandB  *float64 `json:"operand_b,omitempty"`
	Result    float64  `json:"result"`
	CreatedAt string   `json:"created_at"`
}

// HistoryResponse wraps a list of past calculations.
type HistoryResponse struct {
	Calculations []CalculateResponse `json:"calculations"`
}

// ErrorResponse is the exact shape of the JSON we send back when
// something goes wrong — always the same shape, so the frontend can
// handle errors uniformly.
type ErrorResponse struct {
	Error string `json:"error"`
}

// toResponse converts a domain.Calculation (internal model) into a
// CalculateResponse (external API shape). This is the one place where
// the two worlds meet and get translated.
func toResponse(calc domain.Calculation) CalculateResponse {
	return CalculateResponse{
		ID:        calc.ID,
		Operation: string(calc.Operation),
		OperandA:  calc.OperandA,
		OperandB:  calc.OperandB,
		Result:    calc.Result,
		CreatedAt: calc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
