-- Stores every arithmetic operation performed through the API,
-- successful ones only (failed calculations are never persisted —
-- see application.CalculatorService.Calculate).
CREATE TABLE IF NOT EXISTS calculations (
    id          UUID PRIMARY KEY,
    operation   TEXT NOT NULL,
    operand_a   DOUBLE PRECISION NOT NULL,
    operand_b   DOUBLE PRECISION,          -- NULL for unary operations (sqrt)
    result      DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Speeds up "most recent calculations" queries (ORDER BY created_at
-- DESC LIMIT N), which is the only read pattern the app performs.
CREATE INDEX IF NOT EXISTS idx_calculations_created_at
    ON calculations (created_at DESC);