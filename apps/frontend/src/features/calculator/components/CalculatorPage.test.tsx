import { afterEach, describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { CalculatorPage } from "./CalculatorPage";

describe("CalculatorPage", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("shows the Add operation by default", () => {
    render(<CalculatorPage />);
    expect(screen.getByRole("heading", { name: /add/i })).toBeInTheDocument();
  });

  it("switches to another operation's form when its tab is selected", async () => {
    const user = userEvent.setup();
    render(<CalculatorPage />);

    await user.click(screen.getByRole("tab", { name: /square root/i }));

    expect(
      screen.getByRole("heading", { name: /square root/i }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText("Value")).toBeInTheDocument();
  });

  it("completes a full calculate flow after switching tabs", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ result: 3 }),
      }),
    );

    render(<CalculatorPage />);
    await user.click(screen.getByRole("tab", { name: /square root/i }));
    await user.type(screen.getByLabelText("Value"), "9");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(await screen.findByText("= 3")).toBeInTheDocument();
    expect(fetch).toHaveBeenCalledWith(
      "/api/sqrt",
      expect.objectContaining({ body: JSON.stringify({ value: 9 }) }),
    );
  });
});
