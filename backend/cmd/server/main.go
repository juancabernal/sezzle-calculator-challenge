// The server binary: this is the composition root. It's the only
// place in the whole codebase that knows about every concrete
// implementation (in-memory repository, HTTP handlers, the actual
// net/http server) and wires them together.
package main

import (
	"log"
	"net/http"

	httpadapter "github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/http"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/adapters/persistence/memory"
	"github.com/juancabernal/sezzle-calculator-challenge/backend/internal/application"
)

func main() {
	// 1. Build the concrete adapter for persistence. This is the one
	//    line that would change to `postgres.NewRepository(db)` later.
	repo := memory.NewRepository()

	// 2. Inject that adapter into the use case. CalculatorService only
	//    knows it received "something that satisfies CalculationRepository".
	service := application.NewCalculatorService(repo)

	// 3. Inject the use case into the HTTP handler.
	handler := httpadapter.NewHandler(service)

	// 4. Wire routes + CORS.
	router := httpadapter.NewRouter(handler)

	const addr = ":8080"
	log.Printf("calculator backend listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
