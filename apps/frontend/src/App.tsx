import { FormEvent, useState } from "react";
import { calculate, Operator } from "./api/calculator";
import "./App.css";

const OPERATORS: { value: Operator; label: string }[] = [
  { value: "add", label: "+" },
  { value: "subtract", label: "-" },
  { value: "multiply", label: "×" },
  { value: "divide", label: "÷" },
];

export default function App() {
  const [a, setA] = useState("0");
  const [b, setB] = useState("0");
  const [operator, setOperator] = useState<Operator>("add");
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setResult(null);
    try {
      const value = await calculate(Number(a), Number(b), operator);
      setResult(String(value));
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
    }
  }

  return (
    <main className="calculator">
      <h1>Calculator</h1>
      <form onSubmit={handleSubmit}>
        <input
          aria-label="first operand"
          type="number"
          value={a}
          onChange={(e) => setA(e.target.value)}
        />
        <select
          aria-label="operator"
          value={operator}
          onChange={(e) => setOperator(e.target.value as Operator)}
        >
          {OPERATORS.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
        <input
          aria-label="second operand"
          type="number"
          value={b}
          onChange={(e) => setB(e.target.value)}
        />
        <button type="submit">=</button>
      </form>
      <p className={`result ${error ? "error" : ""}`}>
        {error ?? result ?? ""}
      </p>
    </main>
  );
}
