import { describe, expect, it } from "vitest";
import { isUnaryOperation, validateInputs } from "./validation";

describe("isUnaryOperation", () => {
  it("returns true for sqrt", () => {
    expect(isUnaryOperation("sqrt")).toBe(true);
  });

  it("returns false for binary operations like add", () => {
    expect(isUnaryOperation("add")).toBe(false);
  });
});

describe("validateInputs", () => {
  it("rejects an empty first operand", () => {
    expect(validateInputs("add", "", "5")).toBe(
      "Please enter the first number."
    );
  });

  it("rejects a non-numeric first operand", () => {
    expect(validateInputs("add", "abc", "5")).toBe(
      "The first number is not valid."
    );
  });

  it("rejects an empty second operand for a binary operation", () => {
    expect(validateInputs("add", "10", "")).toBe(
      "Please enter the second number."
    );
  });

  it("does not require a second operand for a unary operation", () => {
    expect(validateInputs("sqrt", "9", "")).toBeNull();
  });

  it("rejects division by zero before hitting the network", () => {
    expect(validateInputs("divide", "10", "0")).toBe("Cannot divide by zero.");
  });

  it("accepts valid inputs for a binary operation", () => {
    expect(validateInputs("add", "10", "5")).toBeNull();
  });
});