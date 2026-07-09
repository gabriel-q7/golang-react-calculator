import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, afterEach } from "vitest";
import App from "./App";

describe("App", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("shows the result returned by the API", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ result: 7 }),
      }),
    );

    render(<App />);
    fireEvent.change(screen.getByLabelText("first operand"), {
      target: { value: "3" },
    });
    fireEvent.change(screen.getByLabelText("second operand"), {
      target: { value: "4" },
    });
    fireEvent.click(screen.getByRole("button", { name: "=" }));

    await waitFor(() => expect(screen.getByText("7")).toBeInTheDocument());
  });

  it("shows an error message when the API call fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        json: async () => ({ error: "division by zero" }),
      }),
    );

    render(<App />);
    fireEvent.click(screen.getByRole("button", { name: "=" }));

    await waitFor(() =>
      expect(screen.getByText("division by zero")).toBeInTheDocument(),
    );
  });
});
