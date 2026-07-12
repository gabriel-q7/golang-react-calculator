import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent } from "@/components/ui/card";
import type { Status } from "../engine/reducer";

export interface CalculatorDisplayProps {
  expression: string;
  value: string;
  status: Status;
  errorMessage: string | null;
}

/** The screen of the physical calculator: a running expression line above
 * the current value, with inline error feedback in place of a result. */
export function CalculatorDisplay({
  expression,
  value,
  status,
  errorMessage,
}: CalculatorDisplayProps) {
  return (
    <Card
      className="border-2"
      style={{
        borderColor: "color-mix(in srgb, var(--neon-purple) 45%, transparent)",
        boxShadow:
          "0 0 24px color-mix(in srgb, var(--neon-purple) 25%, transparent), inset 0 0 40px rgba(0,0,0,0.35)",
      }}
    >
      <CardContent className="flex flex-col items-end gap-1 px-6 py-8">
        <p
          className="min-h-[1.5rem] w-full truncate text-right font-mono text-sm text-muted-foreground"
          aria-hidden="true"
        >
          {expression || " "}
        </p>
        <p
          className="w-full truncate text-right font-mono text-4xl font-semibold sm:text-5xl"
          style={{ color: "var(--neon-cyan)", textShadow: "0 0 18px color-mix(in srgb, var(--neon-cyan) 55%, transparent)" }}
          aria-live="polite"
        >
          {status === "loading" ? "…" : value}
        </p>
        {status === "error" && errorMessage ? (
          <Alert variant="destructive" className="mt-2 text-right">
            <AlertDescription>{errorMessage}</AlertDescription>
          </Alert>
        ) : null}
      </CardContent>
    </Card>
  );
}
