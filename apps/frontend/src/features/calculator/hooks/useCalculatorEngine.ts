import { useCallback, useReducer } from "react";
import { calculatorReducer, INITIAL_STATE } from "../engine/reducer";
import { evaluateBinary, evaluateSqrt } from "../engine/operations";
import type { OperatorSymbol } from "../types";

function messageOf(err: unknown): string {
  return err instanceof Error ? err.message : "Calculation failed.";
}

/**
 * Drives the honeycomb keypad: owns the running expression / operand state
 * and talks to the backend (via engine/operations.ts) on operator chaining,
 * "=", and "√". Digit entry, decimals, backspace, and clear are pure local
 * state transitions with no network involved.
 */
export function useCalculatorEngine() {
  const [state, dispatch] = useReducer(calculatorReducer, INITIAL_STATE);

  const inputDigit = useCallback((digit: string) => {
    dispatch({ type: "DIGIT", digit });
  }, []);

  const inputDecimal = useCallback(() => {
    dispatch({ type: "DECIMAL" });
  }, []);

  const backspace = useCallback(() => {
    dispatch({ type: "BACKSPACE" });
  }, []);

  const clear = useCallback(() => {
    dispatch({ type: "CLEAR" });
  }, []);

  const chooseOperator = useCallback(
    async (operator: OperatorSymbol) => {
      const { status, buffer, pendingOperator, hasEnteredOperand, storedValue } = state;
      if (status === "loading") return;

      if (operator === "√") {
        dispatch({ type: "START_EVAL" });
        try {
          const result = await evaluateSqrt(Number(buffer));
          dispatch({ type: "EVAL_SUCCESS_SQRT", result });
        } catch (err) {
          dispatch({ type: "EVAL_ERROR", message: messageOf(err) });
        }
        return;
      }

      if (pendingOperator !== null && hasEnteredOperand) {
        dispatch({ type: "START_EVAL" });
        try {
          const result = await evaluateBinary(pendingOperator, storedValue ?? 0, Number(buffer));
          dispatch({ type: "EVAL_SUCCESS_CHAIN", result, nextOperator: operator });
        } catch (err) {
          dispatch({ type: "EVAL_ERROR", message: messageOf(err) });
        }
        return;
      }

      dispatch({ type: "OPERATOR_REPLACE", operator });
    },
    [state],
  );

  const equals = useCallback(async () => {
    const { status, pendingOperator, storedValue, buffer } = state;
    if (status === "loading" || pendingOperator === null) return;

    dispatch({ type: "START_EVAL" });
    try {
      const result = await evaluateBinary(pendingOperator, storedValue ?? 0, Number(buffer));
      dispatch({ type: "EVAL_SUCCESS_EQUALS", result });
    } catch (err) {
      dispatch({ type: "EVAL_ERROR", message: messageOf(err) });
    }
  }, [state]);

  return {
    display: state.buffer,
    expression: state.expression,
    status: state.status,
    errorMessage: state.errorMessage,
    inputDigit,
    inputDecimal,
    backspace,
    clear,
    chooseOperator,
    equals,
  };
}

export type CalculatorEngine = ReturnType<typeof useCalculatorEngine>;
