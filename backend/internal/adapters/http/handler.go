package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// Handler holds the dependencies the HTTP layer needs — in this case,
// just the use case. It has no idea whether CalculatorService is
// backed by Postgres or an in-memory store; that decision was already
// made and injected when the service was constructed.
type Handler struct {
	service *application.CalculatorService
}

func NewHandler(service *application.CalculatorService) *Handler {
	return &Handler{service: service}
}

// HandleCalculate handles POST /calculate.
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	calc, err := h.service.Calculate(
		r.Context(),
		domain.OperationType(req.Operation),
		req.OperandA,
		req.OperandB,
	)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(calc))
}

// HandleHistory handles GET /history.
func (h *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.History(r.Context(), 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	responses := make([]CalculateResponse, 0, len(history))
	for _, calc := range history {
		responses = append(responses, toResponse(calc))
	}

	writeJSON(w, http.StatusOK, HistoryResponse{Calculations: responses})
}

// writeJSON is a small helper so every handler serializes responses
// the exact same way (correct header, correct status code).
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// writeDomainError maps domain-level errors to the right HTTP status.
// This is the ONLY place in the whole app that decides "this kind of
// business error means 400, not 500" — a deliberate, explicit mapping
// instead of guessing status codes scattered across handlers.
func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrDivisionByZero),
		errors.Is(err, domain.ErrNegativeSqrt),
		errors.Is(err, domain.ErrZeroToNegativePower),
		errors.Is(err, domain.ErrUnknownOperation):
		// These are the client's fault (bad input) -> 400 Bad Request.
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		// Anything else is unexpected -> 500, and we don't leak
		// internal error details to the client.
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
