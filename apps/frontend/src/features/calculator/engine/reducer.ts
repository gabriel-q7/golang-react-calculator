import { formatNumber } from "./format";
import type { BinaryOperatorSymbol } from "../types";

export type Status = "idle" | "loading" | "error";

export interface CalculatorState {
  /** The operand currently being typed. */
  buffer: string;
  /** Operand "A", locked in once an operator is chosen. */
  storedValue: number | null;
  pendingOperator: BinaryOperatorSymbol | null;
  /** Line shown above the main display, e.g. "12 +". */
  expression: string;
  /** Whether the user has typed into `buffer` since the last operator/equals. */
  hasEnteredOperand: boolean;
  status: Status;
  errorMessage: string | null;
}

export const INITIAL_STATE: CalculatorState = {
  buffer: "0",
  storedValue: null,
  pendingOperator: null,
  expression: "",
  hasEnteredOperand: false,
  status: "idle",
  errorMessage: null,
};

const MAX_DIGITS = 15;

export type CalculatorAction =
  | { type: "DIGIT"; digit: string }
  | { type: "DECIMAL" }
  | { type: "BACKSPACE" }
  | { type: "CLEAR" }
  | { type: "OPERATOR_REPLACE"; operator: BinaryOperatorSymbol }
  | { type: "START_EVAL" }
  | { type: "EVAL_SUCCESS_CHAIN"; result: number; nextOperator: BinaryOperatorSymbol }
  | { type: "EVAL_SUCCESS_EQUALS"; result: number }
  | { type: "EVAL_SUCCESS_SQRT"; result: number }
  | { type: "EVAL_ERROR"; message: string };

export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction,
): CalculatorState {
  switch (action.type) {
    case "DIGIT": {
      if (state.status === "loading") return state;
      let buffer: string;
      if (state.buffer === "0") {
        buffer = action.digit === "0" ? "0" : action.digit;
      } else if (state.buffer.replace("-", "").replace(".", "").length >= MAX_DIGITS) {
        buffer = state.buffer;
      } else {
        buffer = state.buffer + action.digit;
      }
      return {
        ...state,
        buffer,
        hasEnteredOperand: true,
        status: "idle",
        errorMessage: null,
      };
    }

    case "DECIMAL": {
      if (state.status === "loading") return state;
      if (state.buffer.includes(".")) {
        return { ...state, status: "idle", errorMessage: null };
      }
      return {
        ...state,
        buffer: state.buffer + ".",
        hasEnteredOperand: true,
        status: "idle",
        errorMessage: null,
      };
    }

    case "BACKSPACE": {
      if (state.status === "loading") return state;
      const trimmed = state.buffer.slice(0, -1);
      return {
        ...state,
        buffer: trimmed === "" || trimmed === "-" ? "0" : trimmed,
        status: "idle",
        errorMessage: null,
      };
    }

    case "CLEAR":
      return { ...INITIAL_STATE };

    case "OPERATOR_REPLACE": {
      if (state.status === "loading") return state;
      if (state.pendingOperator === null) {
        const currentValue = Number(state.buffer);
        return {
          ...state,
          storedValue: currentValue,
          pendingOperator: action.operator,
          expression: `${formatNumber(currentValue)} ${action.operator}`,
          buffer: "0",
          hasEnteredOperand: false,
          status: "idle",
          errorMessage: null,
        };
      }
      return {
        ...state,
        pendingOperator: action.operator,
        expression: `${formatNumber(state.storedValue ?? 0)} ${action.operator}`,
        status: "idle",
        errorMessage: null,
      };
    }

    case "START_EVAL":
      return { ...state, status: "loading", errorMessage: null };

    case "EVAL_SUCCESS_CHAIN":
      return {
        ...state,
        storedValue: action.result,
        pendingOperator: action.nextOperator,
        expression: `${formatNumber(action.result)} ${action.nextOperator}`,
        buffer: "0",
        hasEnteredOperand: false,
        status: "idle",
        errorMessage: null,
      };

    case "EVAL_SUCCESS_EQUALS":
      return {
        ...state,
        expression: `${formatNumber(state.storedValue ?? 0)} ${state.pendingOperator} ${state.buffer} =`,
        buffer: formatNumber(action.result),
        storedValue: null,
        pendingOperator: null,
        hasEnteredOperand: false,
        status: "idle",
        errorMessage: null,
      };

    case "EVAL_SUCCESS_SQRT":
      return {
        ...state,
        expression: `√(${state.buffer})`,
        buffer: formatNumber(action.result),
        hasEnteredOperand: true,
        status: "idle",
        errorMessage: null,
      };

    case "EVAL_ERROR":
      return { ...state, status: "error", errorMessage: action.message };

    default:
      return state;
  }
}
