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

const backspaceKey: KeyDef = {
  label: "⌫",
  ariaLabel: "backspace",
  variant: "control",
  action: { kind: "backspace" },
};
const decimalKey: KeyDef = {
  label: ".",
  ariaLabel: "decimal point",
  variant: "control",
  action: { kind: "decimal" },
};
const clearKey: KeyDef = {
  label: "C",
  ariaLabel: "all clear",
  variant: "control",
  action: { kind: "clear" },
};
const equalsKey: KeyDef = {
  label: "=",
  ariaLabel: "equals",
  variant: "equals",
  action: { kind: "equals" },
};

interface RowDef {
  /** Horizontal start position, in hex-widths, of this row's first key. */
  startOffset: number;
  keys: KeyDef[];
}

/**
 * Conventional physical-calculator arrangement, kept as a single connected
 * honeycomb: digits 7-9/4-6/1-3/0 form a contiguous block on the left, the
 * primary arithmetic operators (÷ × − +) sit in a dedicated column on the
 * right of that block (one per digit row), and the remaining secondary
 * actions (C, √, ^, %, =) form the bottom row. Every key's `startOffset`
 * (row) plus its index (column) is its position in hex-widths.
 *
 * Honeycomb nesting requires each row's `startOffset` to differ from the
 * row immediately above/below it by exactly 0.5 (half a hex-width) — that
 * half-step is what makes one row's hexagons sit in the notches of the
 * next. Every row here alternates 0 / 0.5, including the last one; skip a
 * beat (e.g. two rows in a row at the same offset, or a jump of a full 1)
 * and that row visibly stops lining up with its neighbor.
 */
const ROWS: RowDef[] = [
  { startOffset: 0, keys: [digitKey("7"), digitKey("8"), digitKey("9"), operatorKey("÷", "divide")] },
  { startOffset: 0.5, keys: [digitKey("4"), digitKey("5"), digitKey("6"), operatorKey("×", "multiply")] },
  { startOffset: 0, keys: [digitKey("1"), digitKey("2"), digitKey("3"), operatorKey("−", "subtract")] },
  { startOffset: 0.5, keys: [backspaceKey, digitKey("0"), decimalKey, operatorKey("+", "add")] },
  {
    startOffset: 0,
    keys: [
      clearKey,
      operatorKey("√", "square root", "operator-alt"),
      operatorKey("^", "power", "operator-alt"),
      operatorKey("%", "percent", "operator-alt"),
      equalsKey,
    ],
  },
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
        row.keys.map((key, keyIndex) => (
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
                "--x": row.startOffset + keyIndex,
              } as CSSProperties
            }
          />
        )),
      )}
    </div>
  );
}
