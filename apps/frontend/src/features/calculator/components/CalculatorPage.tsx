import { CalculatorDisplay } from "./CalculatorDisplay";
import { HexKeypad } from "./HexKeypad";
import { useCalculatorEngine } from "../hooks/useCalculatorEngine";

export function CalculatorPage() {
  const engine = useCalculatorEngine();

  return (
    <div className="mx-auto flex w-full max-w-md flex-col gap-6">
      <h1
        className="text-center text-2xl font-bold tracking-[0.2em] uppercase"
        style={{
          color: "var(--neon-pink)",
          textShadow:
            "0 0 12px color-mix(in srgb, var(--neon-pink) 70%, transparent), 0 0 28px color-mix(in srgb, var(--neon-purple) 45%, transparent)",
        }}
      >
        Calculator
      </h1>
      <CalculatorDisplay
        expression={engine.expression}
        value={engine.display}
        status={engine.status}
        errorMessage={engine.errorMessage}
      />
      <HexKeypad engine={engine} />
    </div>
  );
}
