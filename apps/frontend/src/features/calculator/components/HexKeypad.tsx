import type { CSSProperties } from "react";
import { HexButton, type HexKeyVariant } from "./HexButton";
import type { CalculatorEngine } from "../hooks/useCalculatorEngine";
import type { OperatorSymbol } from "../types";
import "./HexKeypad.css";

type KeyAction =
  | { kind: "digit"; value: string }
  | { kind: "decimal" }
  | { kind: "operator"; operator: OperatorSymbol }
  | { kind: "equals" }
  | { kind: "clear" }
  | { kind: "backspace" };

interface KeyDef {
  label: string;
  ariaLabel: string;
  variant: HexKeyVariant;
  action: KeyAction;
}

function digitKey(digit: string): KeyDef {
  return { label: digit, ariaLabel: digit, variant: "digit", action: { kind: "digit", value: digit } };
}

function operatorKey(
  symbol: OperatorSymbol,
  ariaLabel: string,
  variant: HexKeyVariant = "operator",
): KeyDef {
  return { label: symbol, ariaLabel, variant, action: { kind: "operator", operator: symbol } };
}

/**
 * 7 rows x 3 keys, alternating rows staggered by half a hex width (see
 * HexKeypad.css) so every key shares an edge with its neighbors — this is
 * the layout, and its geometry is the layer that makes it a honeycomb
 * rather than a plain grid.
 */
const ROWS: KeyDef[][] = [
  [digitKey("7"), digitKey("8"), digitKey("9")],
  [operatorKey("÷", "divide"), operatorKey("×", "multiply"), operatorKey("−", "subtract")],
  [digitKey("4"), digitKey("5"), digitKey("6")],
  [
    operatorKey("^", "power", "operator-alt"),
    operatorKey("√", "square root", "operator-alt"),
    operatorKey("%", "percent", "operator-alt"),
  ],
  [digitKey("1"), digitKey("2"), digitKey("3")],
  [
    { label: "C", ariaLabel: "all clear", variant: "control", action: { kind: "clear" } },
    { label: "⌫", ariaLabel: "backspace", variant: "control", action: { kind: "backspace" } },
    operatorKey("+", "add"),
  ],
  [
    digitKey("0"),
    { label: ".", ariaLabel: "decimal point", variant: "control", action: { kind: "decimal" } },
    { label: "=", ariaLabel: "equals", variant: "equals", action: { kind: "equals" } },
  ],
];

export interface HexKeypadProps {
  engine: CalculatorEngine;
}

export function HexKeypad({ engine }: HexKeypadProps) {
  const disabled = engine.status === "loading";

  function handleActivate(action: KeyAction) {
    switch (action.kind) {
      case "digit":
        engine.inputDigit(action.value);
        break;
      case "decimal":
        engine.inputDecimal();
        break;
      case "operator":
        void engine.chooseOperator(action.operator);
        break;
      case "equals":
        void engine.equals();
        break;
      case "clear":
        engine.clear();
        break;
      case "backspace":
        engine.backspace();
        break;
    }
  }

  return (
    <div className="hex-keypad" role="group" aria-label="Calculator keypad">
      {ROWS.map((row, rowIndex) =>
        row.map((key, colIndex) => (
          <HexButton
            key={key.ariaLabel}
            label={key.label}
            ariaLabel={key.ariaLabel}
            variant={key.variant}
            disabled={disabled}
            onClick={() => handleActivate(key.action)}
            gridStyle={
              {
                "--row": rowIndex,
                "--col": colIndex,
                "--stagger": rowIndex % 2,
              } as CSSProperties
            }
          />
        )),
      )}
    </div>
  );
}
