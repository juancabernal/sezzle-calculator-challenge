import type { CalculateResponse } from "../../types/calculation";
import styles from "./HistoryList.module.css";

interface HistoryListProps {
  history: CalculateResponse[];
}

// Maps an OperationType to a short symbol, purely for display.
// This is a UI-only concern — it doesn't belong in the domain types.
const OPERATION_SYMBOLS: Record<string, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
  power: "^",
  sqrt: "√",
  percentage: "%",
};

export function HistoryList({ history }: HistoryListProps) {
  if (history.length === 0) {
    return <p className={styles.empty}>No calculations yet.</p>;
  }

  return (
    <ul className={styles.list}>
      {history.map((item) => (
        <li key={item.id} className={styles.item}>
          <span>
            {item.operand_a} {OPERATION_SYMBOLS[item.operation] ?? item.operation}
            {item.operand_b !== undefined ? ` ${item.operand_b}` : ""} ={" "}
            <strong>{item.result}</strong>
          </span>
        </li>
      ))}
    </ul>
  );
}