import { describe, expect, it } from "vitest";
import { formatNumber } from "@/features/calculator/engine/format";

describe("formatNumber", () => {
  it("prints whole numbers with no decimal noise", () => {
    expect(formatNumber(5)).toBe("5");
    expect(formatNumber(-12)).toBe("-12");
    expect(formatNumber(0)).toBe("0");
  });

  it("normalizes negative zero to zero", () => {
    expect(formatNumber(-0)).toBe("0");
  });

  it("hides floating-point noise", () => {
    expect(formatNumber(0.1 + 0.2)).toBe("0.3");
  });

  it("keeps meaningful decimals", () => {
    expect(formatNumber(1 / 3)).toBe("0.333333333333");
  });

  it("renders non-finite values as Error", () => {
    expect(formatNumber(Infinity)).toBe("Error");
    expect(formatNumber(-Infinity)).toBe("Error");
    expect(formatNumber(NaN)).toBe("Error");
  });
});
