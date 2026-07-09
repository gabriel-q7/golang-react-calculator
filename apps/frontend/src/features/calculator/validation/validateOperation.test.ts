import { describe, expect, it } from "vitest";
import { parseFieldValue, validateOperation } from "./validateOperation";

describe("parseFieldValue", () => {
  it("parses a valid integer", () => {
    expect(parseFieldValue("42")).toBe(42);
  });

  it("parses a valid decimal", () => {
    expect(parseFieldValue("3.14")).toBeCloseTo(3.14);
  });

  it("parses a negative number", () => {
    expect(parseFieldValue("-5")).toBe(-5);
  });

  it("returns null for an empty string", () => {
    expect(parseFieldValue("")).toBeNull();
  });

  it("returns null for whitespace-only input", () => {
    expect(parseFieldValue("   ")).toBeNull();
  });

  it("returns null for non-numeric input", () => {
    expect(parseFieldValue("abc")).toBeNull();
  });

  it("returns null for Infinity", () => {
    expect(parseFieldValue("Infinity")).toBeNull();
  });
});

describe("validateOperation", () => {
  it("allows a normal add", () => {
    expect(validateOperation("add", { a: 1, b: 2 })).toBeNull();
  });

  it("rejects division by zero", () => {
    expect(validateOperation("divide", { a: 1, b: 0 })).toBe(
      "Cannot divide by zero.",
    );
  });

  it("allows division by a non-zero divisor", () => {
    expect(validateOperation("divide", { a: 1, b: 2 })).toBeNull();
  });

  it("rejects zero raised to a negative exponent", () => {
    expect(validateOperation("power", { base: 0, exponent: -1 })).toBe(
      "Zero cannot be raised to a negative power.",
    );
  });

  it("rejects a negative base with a fractional exponent", () => {
    expect(validateOperation("power", { base: -4, exponent: 0.5 })).toBe(
      "A negative base requires a whole-number exponent.",
    );
  });

  it("allows a negative base with an integer exponent", () => {
    expect(validateOperation("power", { base: -4, exponent: 2 })).toBeNull();
  });

  it("rejects a negative sqrt operand", () => {
    expect(validateOperation("sqrt", { value: -4 })).toBe(
      "Cannot take the square root of a negative number.",
    );
  });

  it("allows a non-negative sqrt operand", () => {
    expect(validateOperation("sqrt", { value: 9 })).toBeNull();
  });

  it("has no special-case validation for percentage", () => {
    expect(validateOperation("percentage", { value: 200, percent: -10 })).toBeNull();
  });
});
