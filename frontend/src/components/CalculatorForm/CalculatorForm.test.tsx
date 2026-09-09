import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { CalculatorForm } from "./CalculatorForm";

describe("CalculatorForm", () => {
  it("submits with the entered numbers for a binary operation", async () => {
    const onSubmit = vi.fn();
    const user = userEvent.setup();
    render(<CalculatorForm onSubmit={onSubmit} isLoading={false} />);

    // getByLabelText finds the input the same way a screen reader
    // would — by its associated <label>, not by an internal test id.
    // This is the core philosophy of React Testing Library: test what
    // the user actually sees and interacts with.
    await user.type(screen.getByLabelText(/first number/i), "10");
    await user.type(screen.getByLabelText(/second number/i), "4");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(onSubmit).toHaveBeenCalledWith("add", 10, 4);
  });

  it("hides the second input for a unary operation and submits without it", async () => {
    const onSubmit = vi.fn();
    const user = userEvent.setup();
    render(<CalculatorForm onSubmit={onSubmit} isLoading={false} />);

    await user.selectOptions(screen.getByLabelText(/operation/i), "sqrt");

    // The second input should no longer be in the document at all.
    expect(screen.queryByLabelText(/second number/i)).not.toBeInTheDocument();

    await user.type(screen.getByLabelText(/number/i), "9");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(onSubmit).toHaveBeenCalledWith("sqrt", 9, undefined);
  });

  it("shows a validation error and does not call onSubmit when a field is empty", async () => {
    const onSubmit = vi.fn();
    const user = userEvent.setup();
    render(<CalculatorForm onSubmit={onSubmit} isLoading={false} />);

    // Submit with everything empty.
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(screen.getByRole("alert")).toHaveTextContent(
      "Please enter the first number."
    );
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("disables the submit button while loading", () => {
    render(<CalculatorForm onSubmit={vi.fn()} isLoading={true} />);

    expect(screen.getByRole("button", { name: /calculating/i })).toBeDisabled();
  });
});