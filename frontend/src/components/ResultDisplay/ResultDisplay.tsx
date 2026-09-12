import type { CalculateResponse } from "../../types/calculation";
import styles from "./ResultDisplay.module.css";

interface ResultDisplayProps {
  result: CalculateResponse | null;
  error: string | null;
}

// A purely presentational component: given a result or an error, it
// decides what to show. It receives everything via props — no API
// calls, no state of its own.
export function ResultDisplay({ result, error }: ResultDisplayProps) {
  if (error) {
    return (
      <div className={styles.error} role="alert">
        {error}
      </div>
    );
  }

  if (!result) {
    return <div className={styles.placeholder}>Enter numbers to calculate</div>;
  }

  return (
    <div className={styles.result}>
      <span className={styles.label}>Result</span>
      <span className={styles.value}>{result.result}</span>
    </div>
  );
}