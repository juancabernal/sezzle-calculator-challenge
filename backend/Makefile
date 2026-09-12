.PHONY: build run test coverage tidy clean

# Compiles the server binary into bin/server.
build:
	go build -o bin/server ./cmd/server

# Runs the server directly (no separate compile step needed).
run:
	go run ./cmd/server

# Runs all unit tests across every package, verbose output.
test:
	go test ./... -v

# Runs tests with coverage instrumentation, then renders an HTML
# report you can open in a browser to see exactly which lines are
# (and aren't) covered by tests.
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Cleans up dependency declarations (removes unused ones, adds
# missing ones) — the same command we already ran manually earlier.
tidy:
	go mod tidy

# Removes build artifacts.
clean:
	rm -rf bin coverage.out coverage.html