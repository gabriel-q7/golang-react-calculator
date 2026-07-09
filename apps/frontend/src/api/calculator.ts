export type Operator = "add" | "subtract" | "multiply" | "divide";

export interface CalculateResult {
  result: number;
}

export interface CalculateError {
  error: string;
}

export async function calculate(
  a: number,
  b: number,
  operator: Operator,
): Promise<number> {
  const res = await fetch("/api/calculate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ a, b, operator }),
  });

  const body = (await res.json()) as CalculateResult | CalculateError;

  if (!res.ok) {
    throw new Error((body as CalculateError).error ?? "calculation failed");
  }

  return (body as CalculateResult).result;
}
