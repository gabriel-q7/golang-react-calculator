export type OperationId =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "power"
  | "sqrt"
  | "percentage";

export interface FieldConfig {
  /** Key sent in the JSON request body. */
  key: string;
  /** Human-readable field label. */
  label: string;
}

export interface OperationConfig {
  id: OperationId;
  /** Display name, e.g. "Add". */
  label: string;
  /** Short symbol shown next to the label, e.g. "+". */
  symbol: string;
  /** Backend endpoint, e.g. "/api/add". */
  endpoint: string;
  fields: FieldConfig[];
}
