package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/domain"
)

// maxRequestBodyBytes caps how much of a request body we're willing
// to read for POST /calculate. See HandleCalculate for why.
const maxRequestBodyBytes = 4 << 10 // 4 KiB

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
//
//	@Summary		Perform an arithmetic calculation
//	@Description	Executes one operation (add, subtract, multiply, divide, power, sqrt, percentage) and persists the result if successful.
//	@Tags			calculator
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CalculateRequest	true	"Operation and operands"
//	@Success		200		{object}	CalculateResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/calculate [post]
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	// The request is just two floats and an operation name — a few
	// KB is more than generous. Capping it prevents a client from
	// sending an arbitrarily large body just to burn server memory
	// and CPU decoding it.
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

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
//
//	@Summary		List recent calculations
//	@Description	Returns the most recent successful calculations (up to 20), newest first. Failed calculations are never persisted.
//	@Tags			calculator
//	@Produce		json
//	@Success		200	{object}	HistoryResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/history [get]
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
	// The status and headers are already flushed by the time Encode
	// could fail, so there's nothing left to tell the client — but we
	// still log it server-side instead of swallowing it silently, so
	// an unexpected encoding failure doesn't disappear without a trace.
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
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
		errors.Is(err, domain.ErrUnknownOperation),
		errors.Is(err, domain.ErrNonFiniteResult):
		// These are the client's fault (bad input) -> 400 Bad Request.
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		// Anything else is unexpected -> 500, and we don't leak
		// internal error details to the client.
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
