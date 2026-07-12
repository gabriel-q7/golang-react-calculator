import { afterEach, describe, expect, it, vi } from "vitest";
import { evaluateBinary, evaluateSqrt } from "./operations";

function mockFetchOk(result: number) {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ result }),
    }),
  );
}

describe("evaluateBinary", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("posts to the operator's endpoint with the right field names", async () => {
    mockFetchOk(7);
    const result = await evaluateBinary("+", 3, 4);
    expect(result).toBe(7);
    expect(fetch).toHaveBeenCalledWith(
      "/api/add",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ a: 3, b: 4 }),
      }),
    );
  });

  it("maps ^ to /api/power with base/exponent fields", async () => {
    mockFetchOk(1024);
    await evaluateBinary("^", 2, 10);
    expect(fetch).toHaveBeenCalledWith(
      "/api/power",
      expect.objectContaining({ body: JSON.stringify({ base: 2, exponent: 10 }) }),
    );
  });

  it("maps % to /api/percentage with value/percent fields", async () => {
    mockFetchOk(20);
    await evaluateBinary("%", 200, 10);
    expect(fetch).toHaveBeenCalledWith(
      "/api/percentage",
      expect.objectContaining({ body: JSON.stringify({ value: 200, percent: 10 }) }),
    );
  });

  it("rejects divide by zero client-side without calling the API", async () => {
    mockFetchOk(0);
    await expect(evaluateBinary("÷", 1, 0)).rejects.toThrow("Cannot divide by zero.");
    expect(fetch).not.toHaveBeenCalled();
  });
});

describe("evaluateSqrt", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("posts to /api/sqrt with the value field", async () => {
    mockFetchOk(3);
    const result = await evaluateSqrt(9);
    expect(result).toBe(3);
    expect(fetch).toHaveBeenCalledWith(
      "/api/sqrt",
      expect.objectContaining({ body: JSON.stringify({ value: 9 }) }),
    );
  });

  it("rejects negative operands client-side without calling the API", async () => {
    mockFetchOk(0);
    await expect(evaluateSqrt(-4)).rejects.toThrow(
      "Cannot take the square root of a negative number.",
    );
    expect(fetch).not.toHaveBeenCalled();
  });
});
