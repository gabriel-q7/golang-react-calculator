import { afterEach, describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { OperationCard } from "./OperationCard";
import { OPERATIONS } from "../config";

const addConfig = OPERATIONS.find((op) => op.id === "add")!;
const divideConfig = OPERATIONS.find((op) => op.id === "divide")!;
const sqrtConfig = OPERATIONS.find((op) => op.id === "sqrt")!;

describe("OperationCard", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("renders a labeled input per configured field", () => {
    render(<OperationCard config={addConfig} />);
    expect(screen.getByLabelText("A")).toBeInTheDocument();
    expect(screen.getByLabelText("B")).toBeInTheDocument();
  });

  it("submits parsed values and displays the result on success", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ result: 5 }),
      }),
    );

    render(<OperationCard config={addConfig} />);
    await user.type(screen.getByLabelText("A"), "2");
    await user.type(screen.getByLabelText("B"), "3");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(await screen.findByText("= 5")).toBeInTheDocument();
    expect(fetch).toHaveBeenCalledWith(
      "/api/add",
      expect.objectContaining({ body: JSON.stringify({ a: 2, b: 3 }) }),
    );
  });

  it("shows a field error and never calls the API for non-numeric input", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.fn();
    vi.stubGlobal("fetch", fetchMock);

    render(<OperationCard config={addConfig} />);
    await user.type(screen.getByLabelText("A"), "abc");
    await user.type(screen.getByLabelText("B"), "3");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(await screen.findByText("Enter a valid number.")).toBeInTheDocument();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("blocks division by zero client-side without calling the API", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.fn();
    vi.stubGlobal("fetch", fetchMock);

    render(<OperationCard config={divideConfig} />);
    await user.type(screen.getByLabelText("A"), "1");
    await user.type(screen.getByLabelText("B"), "0");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(await screen.findByText("Cannot divide by zero.")).toBeInTheDocument();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("blocks a negative square root client-side without calling the API", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.fn();
    vi.stubGlobal("fetch", fetchMock);

    render(<OperationCard config={sqrtConfig} />);
    await user.type(screen.getByLabelText("Value"), "-4");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    expect(
      await screen.findByText("Cannot take the square root of a negative number."),
    ).toBeInTheDocument();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("surfaces a server-side error returned by the API", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        json: async () => ({ error: "input must be a finite number" }),
      }),
    );

    render(<OperationCard config={addConfig} />);
    await user.type(screen.getByLabelText("A"), "1");
    await user.type(screen.getByLabelText("B"), "2");
    await user.click(screen.getByRole("button", { name: /calculate/i }));

    await waitFor(() =>
      expect(screen.getByText("input must be a finite number")).toBeInTheDocument(),
    );
  });

  it("clears the previous result once an input changes again", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ result: 5 }),
      }),
    );

    render(<OperationCard config={addConfig} />);
    await user.type(screen.getByLabelText("A"), "2");
    await user.type(screen.getByLabelText("B"), "3");
    await user.click(screen.getByRole("button", { name: /calculate/i }));
    expect(await screen.findByText("= 5")).toBeInTheDocument();

    await user.type(screen.getByLabelText("A"), "1");
    expect(screen.queryByText("= 5")).not.toBeInTheDocument();
  });
});
