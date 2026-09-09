import { useEffect } from "react";
import { createHttpCalculatorApi } from "./api/calculatorApi";
import { CalculatorForm } from "./components/CalculatorForm/CalculatorForm";
import { ResultDisplay } from "./components/ResultDisplay/ResultDisplay";
import { HistoryList } from "./components/HistoryList/HistoryList";
import { useCalculator } from "./hooks/useCalculator";
import styles from "./App.module.css";

// This is the composition root: the one place that decides WHICH
// concrete implementation of CalculatorApi to use. Every component
// and hook below this point only knows about the CalculatorApi
// interface, never about fetch or the real backend URL directly.
const calculatorApi = createHttpCalculatorApi(import.meta.env.VITE_API_URL);

function App() {
  const { result, history, error, isLoading, calculate, loadHistory } =
    useCalculator(calculatorApi);

  // Load existing history once, when the app first mounts.
  useEffect(() => {
    loadHistory();
  }, [loadHistory]);

  return (
    <div className={styles.page}>
      <main className={styles.container}>
        <h1 className={styles.title}>Calculator</h1>

        <div className={styles.layout}>
          <section className={styles.panel}>
            <CalculatorForm onSubmit={calculate} isLoading={isLoading} />
            <ResultDisplay result={result} error={error} />
          </section>

          <section className={styles.panel}>
            <h2 className={styles.subtitle}>History</h2>
            <HistoryList history={history} />
          </section>
        </div>
      </main>
    </div>
  );
}

export default App;