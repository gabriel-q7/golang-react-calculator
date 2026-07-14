import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { HexButton } from "@/features/calculator/components/HexButton";

describe("HexButton", () => {
  it("renders its label and accessible name", () => {
    render(
      <HexButton
        label="÷"
        ariaLabel="divide"
        variant="operator"
        onClick={() => {}}
        gridStyle={{}}
      />,
    );
    const button = screen.getByRole("button", { name: "divide" });
    expect(button).toHaveTextContent("÷");
  });

  it("calls onClick when activated", () => {
    const onClick = vi.fn();
    render(
      <HexButton label="7" ariaLabel="7" variant="digit" onClick={onClick} gridStyle={{}} />,
    );
    fireEvent.click(screen.getByRole("button", { name: "7" }));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("does not call onClick when disabled", () => {
    const onClick = vi.fn();
    render(
      <HexButton
        label="="
        ariaLabel="equals"
        variant="equals"
        disabled
        onClick={onClick}
        gridStyle={{}}
      />,
    );
    const button = screen.getByRole("button", { name: "equals" });
    expect(button).toBeDisabled();
    fireEvent.click(button);
    expect(onClick).not.toHaveBeenCalled();
  });

  it("applies the variant class", () => {
    render(
      <HexButton
        label="√"
        ariaLabel="square root"
        variant="operator-alt"
        onClick={() => {}}
        gridStyle={{}}
      />,
    );
    expect(screen.getByRole("button", { name: "square root" })).toHaveClass(
      "hex-btn--operator-alt",
    );
  });
});
