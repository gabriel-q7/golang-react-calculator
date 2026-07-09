import type { OperationConfig } from "./types";

export const OPERATIONS: OperationConfig[] = [
  {
    id: "add",
    label: "Add",
    symbol: "+",
    endpoint: "/api/add",
    fields: [
      { key: "a", label: "A" },
      { key: "b", label: "B" },
    ],
  },
  {
    id: "subtract",
    label: "Subtract",
    symbol: "−",
    endpoint: "/api/subtract",
    fields: [
      { key: "a", label: "A" },
      { key: "b", label: "B" },
    ],
  },
  {
    id: "multiply",
    label: "Multiply",
    symbol: "×",
    endpoint: "/api/multiply",
    fields: [
      { key: "a", label: "A" },
      { key: "b", label: "B" },
    ],
  },
  {
    id: "divide",
    label: "Divide",
    symbol: "÷",
    endpoint: "/api/divide",
    fields: [
      { key: "a", label: "A" },
      { key: "b", label: "B" },
    ],
  },
  {
    id: "power",
    label: "Power",
    symbol: "^",
    endpoint: "/api/power",
    fields: [
      { key: "base", label: "Base" },
      { key: "exponent", label: "Exponent" },
    ],
  },
  {
    id: "sqrt",
    label: "Square Root",
    symbol: "√",
    endpoint: "/api/sqrt",
    fields: [{ key: "value", label: "Value" }],
  },
  {
    id: "percentage",
    label: "Percentage",
    symbol: "%",
    endpoint: "/api/percentage",
    fields: [
      { key: "value", label: "Value" },
      { key: "percent", label: "Percent" },
    ],
  },
];
