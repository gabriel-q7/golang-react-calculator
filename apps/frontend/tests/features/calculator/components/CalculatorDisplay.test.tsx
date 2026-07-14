import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { CalculatorDisplay } from "@/features/calculator/components/CalculatorDisplay";

describe("CalculatorDisplay", () => {
  it("shows the expression and current value", () => {
    render(
      <CalculatorDisplay expression="3 +" value="4" status="idle" errorMessage={null} />,
    );
    expect(screen.getByText("3 +")).toBeInTheDocument();
    expect(screen.getByText("4")).toBeInTheDocument();
  });

  it("shows an ellipsis while loading instead of the stale value", () => {
    render(
      <CalculatorDisplay expression="3 +" value="4" status="loading" errorMessage={null} />,
    );
    expect(screen.getByText("…")).toBeInTheDocument();
  });

  it("shows the error message when status is error", () => {
    render(
      <CalculatorDisplay
        expression="9 ÷"
        value="0"
        status="error"
        errorMessage="Cannot divide by zero."
      />,
    );
    expect(screen.getByRole("alert")).toHaveTextContent("Cannot divide by zero.");
  });

  it("does not render an alert when there is no error", () => {
    render(<CalculatorDisplay expression="" value="0" status="idle" errorMessage={null} />);
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });
});
