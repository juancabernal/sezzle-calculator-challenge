package http

import "net/http"

// NewRouter wires up every route to its handler. Go 1.22+ lets the
// standard http.ServeMux match on method + path directly ("POST /x"),
// so we don't need a third-party router for a small API like this one.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /calculate", h.HandleCalculate)
	mux.HandleFunc("GET /history", h.HandleHistory)

	// Wrap the whole mux with CORS middleware so every route gets it,
	// instead of repeating the header logic in each handler.
	return withCORS(mux)
}

// withCORS wraps a handler so browser-based clients (our React
// frontend, running on a different port during development) are
// allowed to call this API. Without this, the browser blocks the
// response before your frontend code ever sees it — a same-origin
// policy restriction, not a bug in the backend.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Browsers send an OPTIONS "preflight" request before the
		// real POST, just to ask "are you going to allow this?".
		// We answer immediately with 200 and no body — there's
		// nothing else to do for a preflight check.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
