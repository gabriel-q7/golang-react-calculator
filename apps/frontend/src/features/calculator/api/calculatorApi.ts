export interface CalculateSuccess {
  result: number;
}

export interface CalculateFailure {
  error: string;
}

/**
 * Posts operand values to a calculator endpoint and returns the numeric
 * result, or throws an Error with the backend's message on failure.
 */
export async function postCalculate(
  endpoint: string,
  body: Record<string, number>,
): Promise<number> {
  let res: Response;
  try {
    res = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
  } catch {
    throw new Error("Could not reach the server. Please try again.");
  }

  let data: CalculateSuccess | CalculateFailure;
  try {
    data = await res.json();
  } catch {
    throw new Error("Received an unexpected response from the server.");
  }

  if (!res.ok) {
    throw new Error((data as CalculateFailure).error ?? "Calculation failed.");
  }

  return (data as CalculateSuccess).result;
}
