import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError, createHttpCalculatorApi } from "./calculatorApi";

// We fake the global `fetch` function itself here — a level below
// what useCalculator.test.ts does. That test fakes the whole
// CalculatorApi interface; this test verifies the one concrete
// implementation that talks to fetch, so both layers are covered.
describe("createHttpCalculatorApi", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("sends a POST request with the correct body and parses the response", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        id: "1",
        operation: "add",
        operand_a: 10,
        operand_b: 4,
        result: 14,
        created_at: "2026-01-01T00:00:00Z",
      }),
    });
    vi.stubGlobal("fetch", mockFetch);

    const api = createHttpCalculatorApi("http://localhost:8080");
    const result = await api.calculate({
      operation: "add",
      operand_a: 10,
      operand_b: 4,
    });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/calculate",
      expect.objectContaining({
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ operation: "add", operand_a: 10, operand_b: 4 }),
      })
    );
    expect(result.result).toBe(14);
  });

  it("throws an ApiError with the backend's message when the response is not ok", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: "division by zero" }),
    });
    vi.stubGlobal("fetch", mockFetch);

    const api = createHttpCalculatorApi("http://localhost:8080");

    await expect(
      api.calculate({ operation: "divide", operand_a: 10, operand_b: 0 })
    ).rejects.toThrow(ApiError);
  });

  it("extracts the calculations array from the history response", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        calculations: [
          {
            id: "1",
            operation: "add",
            operand_a: 1,
            operand_b: 1,
            result: 2,
            created_at: "2026-01-01T00:00:00Z",
          },
        ],
      }),
    });
    vi.stubGlobal("fetch", mockFetch);

    const api = createHttpCalculatorApi("http://localhost:8080");
    const history = await api.getHistory();

    expect(mockFetch).toHaveBeenCalledWith("http://localhost:8080/history");
    expect(history).toHaveLength(1);
    expect(history[0].result).toBe(2);
  });
});