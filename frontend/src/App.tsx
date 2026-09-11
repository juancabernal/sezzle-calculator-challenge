import { useEffect } from "react";
import { createHttpCalculatorApi } from "./api/calculatorApi";
import { CalculatorForm } from "./components/CalculatorForm/CalculatorForm";
import { ResultDisplay } from "./components/ResultDisplay/ResultDisplay";
import { HistoryList } from "./components/HistoryList/HistoryList";
import { useCalculator } from "./hooks/useCalculator";
import styles from "./App.module.css";

const calculatorApi = createHttpCalculatorApi(import.meta.env.VITE_API_URL);

function App() {
  const { result, history, error, isLoading, calculate, loadHistory } =
    useCalculator(calculatorApi);

  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  return (
    <div className={styles.page}>
      <main className={styles.container}>
        <header className={styles.header}>
          <h1 className={styles.title}>Calculator</h1>
          <p className={styles.subtitle}>
            A small, precise instrument — not a toy.
          </p>
        </header>

        <div className={styles.bento}>
          <section className={`${styles.panel} ${styles.panelMain} glass`}>
            <ResultDisplay result={result} error={error} />
            <CalculatorForm onSubmit={calculate} isLoading={isLoading} />
          </section>

          <section className={`${styles.panel} ${styles.panelHistory} glass`}>
            <h2 className={styles.panelTitle}>History</h2>
            <HistoryList history={history} />
          </section>

          <section className={`${styles.panel} ${styles.panelStack} glass`}>
            <span className={styles.stackLabel}>Built with</span>
            <span className={styles.stackValue}>
              Go · React · TypeScript
            </span>
            <span className={styles.stackSub}>Hexagonal architecture</span>
          </section>
        </div>
      </main>
    </div>
  );
}

export default App;