package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpadapter "github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/http"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/memory"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
)

// newTestRouter builds a fresh router with an isolated in-memory
// repository for each test, so tests never leak state into each other.
func newTestRouter() http.Handler {
	repo := memory.NewRepository()
	service := application.NewCalculatorService(repo)
	handler := httpadapter.NewHandler(service)
	return httpadapter.NewRouter(handler)
}

func TestHandleCalculate_Success(t *testing.T) {
	router := newTestRouter()

	body := strings.NewReader(`{"operation":"multiply","operand_a":6,"operand_b":7}`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", body)
	req.Header.Set("Content-Type", "application/json")

	// httptest.NewRecorder acts as a fake ResponseWriter: it captures
	// whatever the handler writes (status code, headers, body) so we
	// can inspect it afterwards — no real network socket involved.
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp httpadapter.CalculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 42 {
		t.Errorf("Result = %v, want 42", resp.Result)
	}
}

func TestHandleCalculate_DivisionByZero_ReturnsBadRequest(t *testing.T) {
	router := newTestRouter()

	body := strings.NewReader(`{"operation":"divide","operand_a":5,"operand_b":0}`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleCalculate_NonFiniteResult_ReturnsBadRequest(t *testing.T) {
	router := newTestRouter()

	// (-8)^0.5 is not a real number -> math.Pow returns NaN, which
	// encoding/json cannot marshal. This must surface as a clean 400
	// with a JSON error body, never as a 200 with an empty/broken body.
	body := strings.NewReader(`{"operation":"power","operand_a":-8,"operand_b":0.5}`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d. Body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}

	var resp httpadapter.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected a non-empty error message")
	}
}

func TestHandleCalculate_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	router := newTestRouter()

	body := strings.NewReader(`not valid json`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleHistory_ReturnsOnlySuccessfulCalculations(t *testing.T) {
	router := newTestRouter()

	// One successful calculation...
	successBody := strings.NewReader(`{"operation":"add","operand_a":1,"operand_b":1}`)
	successReq := httptest.NewRequest(http.MethodPost, "/calculate", successBody)
	router.ServeHTTP(httptest.NewRecorder(), successReq)

	// ...and one that fails.
	failBody := strings.NewReader(`{"operation":"divide","operand_a":1,"operand_b":0}`)
	failReq := httptest.NewRequest(http.MethodPost, "/calculate", failBody)
	router.ServeHTTP(httptest.NewRecorder(), failReq)

	historyReq := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, historyReq)

	var resp httpadapter.HistoryResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Calculations) != 1 {
		t.Errorf("Calculations length = %d, want 1 (failed calc should be excluded)", len(resp.Calculations))
	}
}

func TestCORSHeaders_ArePresent(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("preflight status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
}
