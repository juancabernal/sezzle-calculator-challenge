import { useState } from "react";
import type { OperationType } from "../../types/calculation";
import { isUnaryOperation, validateInputs } from "../../utils/validation";
import styles from "./CalculatorForm.module.css";

interface CalculatorFormProps {
  // The form doesn't know HOW the calculation happens (API call,
  // logging, anything) — it just reports "the user wants to
  // calculate this" upward, via a callback the parent provides.
  onSubmit: (
    operation: OperationType,
    operandA: number,
    operandB?: number
  ) => void;
  isLoading: boolean;
}

const OPERATIONS: { value: OperationType; label: string }[] = [
  { value: "add", label: "Addition (+)" },
  { value: "subtract", label: "Subtraction (−)" },
  { value: "multiply", label: "Multiplication (×)" },
  { value: "divide", label: "Division (÷)" },
  { value: "power", label: "Power (^)" },
  { value: "sqrt", label: "Square root (√)" },
  { value: "percentage", label: "Percentage (%)" },
];

export function CalculatorForm({ onSubmit, isLoading }: CalculatorFormProps) {
  const [operation, setOperation] = useState<OperationType>("add");
  const [operandA, setOperandA] = useState("");
  const [operandB, setOperandB] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);

  const unary = isUnaryOperation(operation);

  function handleSubmit(event: React.FormEvent) {
    event.preventDefault();

    const error = validateInputs(operation, operandA, unary ? "0" : operandB);
    if (error) {
      setValidationError(error);
      return;
    }

    setValidationError(null);
    onSubmit(
      operation,
      Number(operandA),
      unary ? undefined : Number(operandB)
    );
  }

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <label className={styles.field}>
        Operation
        <select
          value={operation}
          onChange={(e) => setOperation(e.target.value as OperationType)}
        >
          {OPERATIONS.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
      </label>

      <label className={styles.field}>
        {unary ? "Number" : "First number"}
        <input
          type="text"
          inputMode="decimal"
          value={operandA}
          onChange={(e) => setOperandA(e.target.value)}
          placeholder="0"
        />
      </label>

      {!unary && (
        <label className={styles.field}>
          Second number
          <input
            type="text"
            inputMode="decimal"
            value={operandB}
            onChange={(e) => setOperandB(e.target.value)}
            placeholder="0"
          />
        </label>
      )}

      {validationError && (
        <p className={styles.error} role="alert">
          {validationError}
        </p>
      )}

      <button type="submit" disabled={isLoading}>
        {isLoading ? "Calculating..." : "Calculate"}
      </button>
    </form>
  );
}