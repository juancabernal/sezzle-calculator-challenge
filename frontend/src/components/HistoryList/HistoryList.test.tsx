import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { HistoryList } from "./HistoryList";
import type { CalculateResponse } from "../../types/calculation";

describe("HistoryList", () => {
  it("shows an empty state when there is no history", () => {
    render(<HistoryList history={[]} />);
    expect(screen.getByText(/no calculations yet/i)).toBeInTheDocument();
  });

  it("renders one entry per calculation, using the correct operation symbol", () => {
    const history: CalculateResponse[] = [
      {
        id: "1",
        operation: "add",
        operand_a: 10,
        operand_b: 4,
        result: 14,
        created_at: "2026-01-01T00:00:00Z",
      },
      {
        id: "2",
        operation: "sqrt",
        operand_a: 9,
        result: 3,
        created_at: "2026-01-01T00:00:01Z",
      },
    ];

    render(<HistoryList history={history} />);

    expect(screen.getByText(/10 \+ 4/)).toBeInTheDocument();
    expect(screen.getByText(/9 √/)).toBeInTheDocument();
  });
});