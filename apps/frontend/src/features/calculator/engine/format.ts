/**
 * Formats a number for the calculator display: whole numbers print with no
 * decimal noise, everything else is rounded to 12 significant digits to
 * hide binary floating-point artifacts (e.g. 0.1 + 0.2).
 */
export function formatNumber(value: number): string {
  if (!Number.isFinite(value)) {
    return "Error";
  }
  const normalized = Object.is(value, -0) ? 0 : value;
  if (Number.isInteger(normalized) && Math.abs(normalized) < 1e15) {
    return normalized.toString();
  }
  return Number(normalized.toPrecision(12)).toString();
}
