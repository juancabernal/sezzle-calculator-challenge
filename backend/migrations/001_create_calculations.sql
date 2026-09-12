-- Stores every arithmetic operation performed through the API,
-- successful ones only (failed calculations are never persisted —
-- see application.CalculatorService.Calculate).
CREATE TABLE IF NOT EXISTS calculations (
    -- BIGSERIAL is monotonically increasing and assigned atomically by
    -- Postgres on insert, unlike created_at: two requests handled by
    -- concurrent goroutines can be timestamped in the very same
    -- microsecond (TIMESTAMPTZ's own resolution), which would make
    -- "most recent first" ties break arbitrarily. seq exists purely
    -- as a reliable ordering key — the domain layer never reads it.
    seq         BIGSERIAL PRIMARY KEY,
    id          UUID NOT NULL UNIQUE,
    operation   TEXT NOT NULL,
    operand_a   DOUBLE PRECISION NOT NULL,
    operand_b   DOUBLE PRECISION,          -- NULL for unary operations (sqrt)
    result      DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Speeds up "most recent calculations" queries (ORDER BY seq DESC
-- LIMIT N), which is the only read pattern the app performs.
CREATE INDEX IF NOT EXISTS idx_calculations_seq
    ON calculations (seq DESC);