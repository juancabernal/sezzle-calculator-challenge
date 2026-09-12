import { describe, expect, it, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useCalculator } from "./useCalculator";
import { ApiError, type CalculatorApi } from "../api/calculatorApi";
import type { CalculateResponse } from "../types/calculation";

// A fake, in-memory implementation of the CalculatorApi port — the
// same idea as the memory.Repository we wrote in Go. Because
// useCalculator depends on the CalculatorApi interface (not on fetch
// directly), we can substitute this fake with zero mocking libraries.
function createFakeApi(overrides?: Partial<CalculatorApi>): CalculatorApi {
  return {
    async calculate(request) {
      return {
        id: "fake-id",
        operation: request.operation,
        operand_a: request.operand_a,
        operand_b: request.operand_b,
        result: request.operand_a + (request.operand_b ?? 0),
        created_at: new Date().toISOString(),
      };
    },
    async getHistory() {
      return [];
    },
    ...overrides,
  };
}

describe("useCalculator", () => {
  it("stores the result after a successful calculation", async () => {
    const api = createFakeApi();
    const { result } = renderHook(() => useCalculator(api));

    await act(async () => {
      await result.current.calculate("add", 10, 4);
    });

    expect(result.current.result?.result).toBe(14);
    expect(result.current.error).toBeNull();
  });

  it("prepends successful calculations to history", async () => {
    const api = createFakeApi();
    const { result } = renderHook(() => useCalculator(api));

    await act(async () => {
      await result.current.calculate("add", 1, 1);
    });
    await act(async () => {
      await result.current.calculate("add", 2, 2);
    });

    expect(result.current.history).toHaveLength(2);
    // Most recent calculation should be first.
    expect(result.current.history[0].result).toBe(4);
  });

  it("surfaces the exact message from an ApiError", async () => {
    const api = createFakeApi({
      calculate: vi.fn().mockRejectedValue(new ApiError("division by zero", 400)),
    });
    const { result } = renderHook(() => useCalculator(api));

    await act(async () => {
      await result.current.calculate("divide", 10, 0);
    });

    expect(result.current.error).toBe("division by zero");
    expect(result.current.result).toBeNull();
  });

  it("shows a generic message for non-ApiError failures", async () => {
    const api = createFakeApi({
      calculate: vi.fn().mockRejectedValue(new Error("network down")),
    });
    const { result } = renderHook(() => useCalculator(api));

    await act(async () => {
      await result.current.calculate("add", 1, 1);
    });

    expect(result.current.error).toBe("Could not reach the calculator service.");
  });

  it("loads history on demand", async () => {
    const fakeHistory: CalculateResponse[] = [
      {
        id: "1",
        operation: "add",
        operand_a: 1,
        operand_b: 1,
        result: 2,
        created_at: new Date().toISOString(),
      },
    ];
    const api = createFakeApi({ getHistory: async () => fakeHistory });
    const { result } = renderHook(() => useCalculator(api));

    await act(async () => {
      await result.current.loadHistory();
    });

    await waitFor(() => {
      expect(result.current.history).toEqual(fakeHistory);
    });
  });
});