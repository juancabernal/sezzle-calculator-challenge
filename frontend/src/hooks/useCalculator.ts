import { useCallback, useState } from "react";
import type { CalculatorApi } from "../api/calculatorApi";
import { ApiError } from "../api/calculatorApi";
import type { CalculateResponse, OperationType } from "../types/calculation";

// useCalculator is the use case of the frontend: it orchestrates calls
// to the CalculatorApi port and exposes just the state a component
// needs to render. Components never call the API directly — they
// call the functions this hook returns.
//
// Notice the hook takes `api: CalculatorApi` (the interface) as a
// parameter, not a hardcoded import of the real fetch-based client.
// That's what lets tests pass in a fake implementation instead.
export function useCalculator(api: CalculatorApi) {
  const [result, setResult] = useState<CalculateResponse | null>(null);
  const [history, setHistory] = useState<CalculateResponse[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const calculate = useCallback(
    async (operation: OperationType, operandA: number, operandB?: number) => {
      setIsLoading(true);
      setError(null);

      try {
        const response = await api.calculate({
          operation,
          operand_a: operandA,
          operand_b: operandB,
        });
        setResult(response);
        setHistory((prev) => [response, ...prev]);
      } catch (err) {
        // ApiError carries the exact message the Go backend sent
        // (e.g. "division by zero") — we show that directly instead
        // of a generic "something went wrong".
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError("Could not reach the calculator service.");
        }
        setResult(null);
      } finally {
        setIsLoading(false);
      }
    },
    [api]
  );

  const loadHistory = useCallback(async () => {
    try {
      const items = await api.getHistory();
      setHistory(items);
    } catch {
      // Silently ignore: history is a "nice to have" on load, not
      // worth showing an error banner for.
    }
  }, [api]);

  return { result, history, error, isLoading, calculate, loadHistory };
}