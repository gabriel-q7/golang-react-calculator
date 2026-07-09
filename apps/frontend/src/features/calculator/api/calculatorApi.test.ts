import { afterEach, describe, expect, it, vi } from "vitest";
import { postCalculate } from "./calculatorApi";

describe("postCalculate", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("posts the body to the given endpoint and returns the result", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ result: 42 }),
    });
    vi.stubGlobal("fetch", fetchMock);

    const result = await postCalculate("/api/add", { a: 40, b: 2 });

    expect(result).toBe(42);
    expect(fetchMock).toHaveBeenCalledWith(
      "/api/add",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ a: 40, b: 2 }),
      }),
    );
  });

  it("throws the backend error message on a non-OK response", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        json: async () => ({ error: "division by zero" }),
      }),
    );

    await expect(postCalculate("/api/divide", { a: 1, b: 0 })).rejects.toThrow(
      "division by zero",
    );
  });

  it("throws a generic error when the network request fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("network down")),
    );

    await expect(postCalculate("/api/add", { a: 1, b: 2 })).rejects.toThrow(
      "Could not reach the server",
    );
  });

  it("throws a generic error when the response body isn't valid JSON", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => {
          throw new SyntaxError("Unexpected token");
        },
      }),
    );

    await expect(postCalculate("/api/add", { a: 1, b: 2 })).rejects.toThrow(
      "unexpected response",
    );
  });
});
