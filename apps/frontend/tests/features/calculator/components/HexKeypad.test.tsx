import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { HexKeypad } from "@/features/calculator/components/HexKeypad";
import type { CalculatorEngine } from "@/features/calculator/hooks/useCalculatorEngine";

function makeEngine(overrides: Partial<CalculatorEngine> = {}): CalculatorEngine {
  return {
    display: "0",
    expression: "",
    status: "idle",
    errorMessage: null,
    inputDigit: vi.fn(),
    inputDecimal: vi.fn(),
    backspace: vi.fn(),
    clear: vi.fn(),
    chooseOperator: vi.fn().mockResolvedValue(undefined),
    equals: vi.fn().mockResolvedValue(undefined),
    ...overrides,
  };
}

describe("HexKeypad", () => {
  it("renders all 21 keys as an accessible group", () => {
    render(<HexKeypad engine={makeEngine()} />);
    expect(screen.getByRole("group", { name: "Calculator keypad" })).toBeInTheDocument();
    expect(screen.getAllByRole("button")).toHaveLength(21);
  });

  it("routes digit keys to engine.inputDigit", () => {
    const engine = makeEngine();
    render(<HexKeypad engine={engine} />);
    fireEvent.click(screen.getByRole("button", { name: "7" }));
    expect(engine.inputDigit).toHaveBeenCalledWith("7");
  });

  it("routes the decimal key to engine.inputDecimal", () => {
    const engine = makeEngine();
    render(<HexKeypad engine={engine} />);
    fireEvent.click(screen.getByRole("button", { name: "decimal point" }));
    expect(engine.inputDecimal).toHaveBeenCalledTimes(1);
  });

  it("routes operator keys to engine.chooseOperator, including √", () => {
    const engine = makeEngine();
    render(<HexKeypad engine={engine} />);
    fireEvent.click(screen.getByRole("button", { name: "divide" }));
    expect(engine.chooseOperator).toHaveBeenCalledWith("÷");
    fireEvent.click(screen.getByRole("button", { name: "square root" }));
    expect(engine.chooseOperator).toHaveBeenCalledWith("√");
  });

  it("routes = to engine.equals and C/⌫ to clear/backspace", () => {
    const engine = makeEngine();
    render(<HexKeypad engine={engine} />);
    fireEvent.click(screen.getByRole("button", { name: "equals" }));
    expect(engine.equals).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByRole("button", { name: "all clear" }));
    expect(engine.clear).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByRole("button", { name: "backspace" }));
    expect(engine.backspace).toHaveBeenCalledTimes(1);
  });

  it("disables every key while a request is loading", () => {
    const engine = makeEngine({ status: "loading" });
    render(<HexKeypad engine={engine} />);
    for (const button of screen.getAllByRole("button")) {
      expect(button).toBeDisabled();
    }
  });
});
