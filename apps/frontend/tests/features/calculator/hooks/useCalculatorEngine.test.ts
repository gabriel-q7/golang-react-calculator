import { act, renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { useCalculatorEngine } from "@/features/calculator/hooks/useCalculatorEngine";

function mockFetchOnce(result: number) {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ result }),
    }),
  );
}

function mockFetchError(message: string) {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: false,
      json: async () => ({ error: message }),
    }),
  );
}

describe("useCalculatorEngine", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("builds up a multi-digit number in the display", () => {
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => {
      result.current.inputDigit("1");
      result.current.inputDigit("2");
      result.current.inputDecimal();
      result.current.inputDigit("5");
    });
    expect(result.current.display).toBe("12.5");
  });

  it("ignores a second decimal point", () => {
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => {
      result.current.inputDigit("1");
      result.current.inputDecimal();
      result.current.inputDigit("5");
      result.current.inputDecimal();
      result.current.inputDigit("3");
    });
    expect(result.current.display).toBe("1.53");
  });

  it("backspaces one character at a time, bottoming out at 0", () => {
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => {
      result.current.inputDigit("1");
      result.current.inputDigit("2");
    });
    act(() => result.current.backspace());
    expect(result.current.display).toBe("1");
    act(() => result.current.backspace());
    expect(result.current.display).toBe("0");
    act(() => result.current.backspace());
    expect(result.current.display).toBe("0");
  });

  it("performs add via the backend on operator + digits + equals", async () => {
    mockFetchOnce(7);
    const { result } = renderHook(() => useCalculatorEngine());

    act(() => result.current.inputDigit("3"));
    await act(async () => {
      await result.current.chooseOperator("+");
    });
    expect(result.current.expression).toBe("3 +");
    expect(result.current.display).toBe("0");

    act(() => result.current.inputDigit("4"));
    await act(async () => {
      await result.current.equals();
    });

    await waitFor(() => expect(result.current.display).toBe("7"));
    expect(result.current.expression).toBe("3 + 4 =");
    expect(fetch).toHaveBeenCalledWith(
      "/api/add",
      expect.objectContaining({ body: JSON.stringify({ a: 3, b: 4 }) }),
    );
  });

  it("chains operators: pressing an operator mid-entry evaluates the pending one first", async () => {
    mockFetchOnce(7); // 3 + 4 = 7
    const { result } = renderHook(() => useCalculatorEngine());

    act(() => result.current.inputDigit("3"));
    await act(async () => {
      await result.current.chooseOperator("+");
    });
    act(() => result.current.inputDigit("4"));

    await act(async () => {
      await result.current.chooseOperator("÷");
    });

    await waitFor(() => expect(result.current.expression).toBe("7 ÷"));
    expect(result.current.display).toBe("0");
    expect(fetch).toHaveBeenLastCalledWith(
      "/api/add",
      expect.objectContaining({ body: JSON.stringify({ a: 3, b: 4 }) }),
    );
  });

  it("replaces the pending operator when pressed again with no new digits", async () => {
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => result.current.inputDigit("5"));
    await act(async () => {
      await result.current.chooseOperator("+");
    });
    await act(async () => {
      await result.current.chooseOperator("×");
    });
    expect(result.current.expression).toBe("5 ×");
  });

  it("applies √ immediately to the current display without disturbing a pending operator", async () => {
    mockFetchOnce(4); // √16 = 4
    const { result } = renderHook(() => useCalculatorEngine());

    act(() => result.current.inputDigit("2"));
    await act(async () => {
      await result.current.chooseOperator("×"); // stored=2, pending=×
    });
    act(() => {
      result.current.inputDigit("1");
      result.current.inputDigit("6");
    });

    await act(async () => {
      await result.current.chooseOperator("√"); // √16 = 4, replaces buffer only
    });
    await waitFor(() => expect(result.current.display).toBe("4"));

    mockFetchOnce(8); // 2 × 4 = 8
    await act(async () => {
      await result.current.equals();
    });
    await waitFor(() => expect(result.current.display).toBe("8"));
    expect(fetch).toHaveBeenLastCalledWith(
      "/api/multiply",
      expect.objectContaining({ body: JSON.stringify({ a: 2, b: 4 }) }),
    );
  });

  it("blocks divide by zero client-side and shows the error without resetting the buffer", async () => {
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => result.current.inputDigit("9"));
    await act(async () => {
      await result.current.chooseOperator("÷");
    });
    act(() => result.current.inputDigit("0"));
    await act(async () => {
      await result.current.equals();
    });

    expect(result.current.status).toBe("error");
    expect(result.current.errorMessage).toBe("Cannot divide by zero.");
    expect(result.current.display).toBe("0");
  });

  it("surfaces a backend error message and clears it on the next digit press", async () => {
    mockFetchError("input must be a finite number");
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => result.current.inputDigit("1"));
    await act(async () => {
      await result.current.chooseOperator("+");
    });
    act(() => result.current.inputDigit("2"));
    await act(async () => {
      await result.current.equals();
    });

    await waitFor(() => expect(result.current.status).toBe("error"));
    expect(result.current.errorMessage).toBe("input must be a finite number");

    act(() => result.current.inputDigit("5"));
    expect(result.current.status).toBe("idle");
    expect(result.current.errorMessage).toBeNull();
  });

  it("clear resets to the initial state", async () => {
    mockFetchOnce(7);
    const { result } = renderHook(() => useCalculatorEngine());
    act(() => result.current.inputDigit("3"));
    await act(async () => {
      await result.current.chooseOperator("+");
    });
    act(() => result.current.inputDigit("4"));
    act(() => result.current.clear());

    expect(result.current.display).toBe("0");
    expect(result.current.expression).toBe("");
    expect(result.current.status).toBe("idle");
  });
});
