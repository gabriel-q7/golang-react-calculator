import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { CalculatorPage } from "./CalculatorPage";

function mockFetchOk(result: number) {
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

describe("CalculatorPage", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("computes 3 + 4 via the backend and displays the result", async () => {
    mockFetchOk(7);
    render(<CalculatorPage />);

    fireEvent.click(screen.getByRole("button", { name: "3" }));
    fireEvent.click(screen.getByRole("button", { name: "add" }));
    fireEvent.click(screen.getByRole("button", { name: "4" }));
    fireEvent.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() =>
      expect(screen.getByText("7", { selector: "p" })).toBeInTheDocument(),
    );
    expect(fetch).toHaveBeenCalledWith(
      "/api/add",
      expect.objectContaining({ body: JSON.stringify({ a: 3, b: 4 }) }),
    );
  });

  it("shows a backend error without crashing the keypad", async () => {
    mockFetchError("division by zero");
    render(<CalculatorPage />);

    fireEvent.click(screen.getByRole("button", { name: "9" }));
    fireEvent.click(screen.getByRole("button", { name: "divide" }));
    fireEvent.click(screen.getByRole("button", { name: "0" }));
    fireEvent.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
  });

  it("blocks a client-side-invalid divide by zero before calling the API", async () => {
    render(<CalculatorPage />);

    fireEvent.click(screen.getByRole("button", { name: "9" }));
    fireEvent.click(screen.getByRole("button", { name: "divide" }));
    fireEvent.click(screen.getByRole("button", { name: "0" }));
    fireEvent.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() =>
      expect(screen.getByRole("alert")).toHaveTextContent("Cannot divide by zero."),
    );
  });

  it("computes √9 immediately without pressing equals", async () => {
    mockFetchOk(3);
    render(<CalculatorPage />);

    fireEvent.click(screen.getByRole("button", { name: "9" }));
    fireEvent.click(screen.getByRole("button", { name: "square root" }));

    await waitFor(() =>
      expect(screen.getByText("3", { selector: "p" })).toBeInTheDocument(),
    );
    expect(fetch).toHaveBeenCalledWith(
      "/api/sqrt",
      expect.objectContaining({ body: JSON.stringify({ value: 9 }) }),
    );
  });

  it("clears back to 0 on all clear", () => {
    render(<CalculatorPage />);
    fireEvent.click(screen.getByRole("button", { name: "5" }));
    fireEvent.click(screen.getByRole("button", { name: "all clear" }));
    expect(screen.getByText("0", { selector: "p" })).toBeInTheDocument();
  });
});
