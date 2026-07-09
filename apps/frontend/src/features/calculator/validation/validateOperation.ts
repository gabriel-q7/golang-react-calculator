import type { OperationId } from "../types";

/**
 * Parses a raw form field value into a finite number, or null if the
 * input is empty, not a number, or NaN/Infinity.
 */
export function parseFieldValue(raw: string): number | null {
  if (raw.trim() === "") {
    return null;
  }
  const value = Number(raw);
  return Number.isFinite(value) ? value : null;
}

/**
 * Domain-level validation mirroring the backend's rules, run client-side
 * so obviously-invalid requests never round-trip to the API. The backend
 * remains the source of truth and re-validates independently.
 */
export function validateOperation(
  id: OperationId,
  values: Record<string, number>,
): string | null {
  switch (id) {
    case "divide":
      if (values.b === 0) {
        return "Cannot divide by zero.";
      }
      return null;

    case "power":
      if (values.base === 0 && values.exponent < 0) {
        return "Zero cannot be raised to a negative power.";
      }
      if (values.base < 0 && !Number.isInteger(values.exponent)) {
        return "A negative base requires a whole-number exponent.";
      }
      return null;

    case "sqrt":
      if (values.value < 0) {
        return "Cannot take the square root of a negative number.";
      }
      return null;

    default:
      return null;
  }
}
