import type {
  CalculateRequest,
  CalculateResponse,
  ApiErrorResponse,
} from "../types/calculation";

// CalculatorApi is the port: the contract that anything talking to the
// backend must satisfy. Components and hooks depend on THIS interface,
// never on `fetch` directly. This is what lets tests substitute a fake
// implementation without mocking the network.
export interface CalculatorApi {
  calculate(request: CalculateRequest): Promise<CalculateResponse>;
  getHistory(): Promise<CalculateResponse[]>;
}

// A typed error so callers (the hook, eventually the UI) can
// distinguish "the backend rejected this calculation" from other
// failures, and show the exact message the backend provided.
export class ApiError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

// createHttpCalculatorApi is the concrete adapter: it implements
// CalculatorApi using fetch against a real backend. baseUrl is
// injected as a parameter (not hardcoded), so this same function
// works against localhost during development and against a deployed
// URL in production — the decision of WHICH url lives at the call
// site (main.tsx), not buried inside this file.
export function createHttpCalculatorApi(baseUrl: string): CalculatorApi {
  async function handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const body: ApiErrorResponse = await response.json();
      throw new ApiError(body.error, response.status);
    }
    return response.json() as Promise<T>;
  }

  return {
    async calculate(request) {
      const response = await fetch(`${baseUrl}/calculate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(request),
      });
      return handleResponse<CalculateResponse>(response);
    },

    async getHistory() {
      const response = await fetch(`${baseUrl}/history`);
      const data = await handleResponse<{ calculations: CalculateResponse[] }>(
        response
      );
      return data.calculations;
    },
  };
}