import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { ResultDisplay } from "./ResultDisplay";
import type { CalculateResponse } from "../../types/calculation";

const sampleResult: CalculateResponse = {
  id: "1",
  operation: "add",
  operand_a: 10,
  operand_b: 4,
  result: 14,
  created_at: "2026-01-01T00:00:00Z",
};

describe("ResultDisplay", () => {
  it("shows a placeholder when there is no result or error yet", () => {
    render(<ResultDisplay result={null} error={null} />);
    expect(screen.getByText(/enter numbers to calculate/i)).toBeInTheDocument();
  });

  it("shows the result value when a calculation succeeded", () => {
    render(<ResultDisplay result={sampleResult} error={null} />);
    expect(screen.getByText("14")).toBeInTheDocument();
  });

  it("prioritizes showing the error over any stale result", () => {
    render(<ResultDisplay result={sampleResult} error="division by zero" />);
    expect(screen.getByRole("alert")).toHaveTextContent("division by zero");
    // The old result must not still be shown alongside the error.
    expect(screen.queryByText("14")).not.toBeInTheDocument();
  });
});