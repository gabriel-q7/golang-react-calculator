import { OPERATIONS } from "../config";
import { postCalculate } from "../api/calculatorApi";
import { validateOperation } from "../validation/validateOperation";
import type { BinaryOperatorSymbol, OperationConfig, OperationId } from "../types";

const OPERATIONS_BY_ID = Object.fromEntries(
  OPERATIONS.map((operation) => [operation.id, operation]),
) as Record<OperationId, OperationConfig>;

const BINARY_OPERATOR_TO_OPERATION_ID: Record<BinaryOperatorSymbol, OperationId> = {
  "+": "add",
  "−": "subtract",
  "×": "multiply",
  "÷": "divide",
  "^": "power",
  "%": "percentage",
};

function buildValues(
  config: OperationConfig,
  a: number,
  b?: number,
): Record<string, number> {
  const values: Record<string, number> = { [config.fields[0].key]: a };
  if (config.fields[1] && b !== undefined) {
    values[config.fields[1].key] = b;
  }
  return values;
}

/**
 * Runs a binary keypad operator (everything but √) against the backend,
 * pre-validating client-side first so obviously-invalid requests (e.g.
 * divide by zero) never round-trip to the API.
 */
export async function evaluateBinary(
  operator: BinaryOperatorSymbol,
  a: number,
  b: number,
): Promise<number> {
  const id = BINARY_OPERATOR_TO_OPERATION_ID[operator];
  const config = OPERATIONS_BY_ID[id];
  const values = buildValues(config, a, b);

  const domainError = validateOperation(id, values);
  if (domainError) {
    throw new Error(domainError);
  }

  return postCalculate(config.endpoint, values);
}

/** Runs √ (the only unary keypad operator) against the backend. */
export async function evaluateSqrt(value: number): Promise<number> {
  const config = OPERATIONS_BY_ID.sqrt;
  const values = buildValues(config, value);

  const domainError = validateOperation("sqrt", values);
  if (domainError) {
    throw new Error(domainError);
  }

  return postCalculate(config.endpoint, values);
}
