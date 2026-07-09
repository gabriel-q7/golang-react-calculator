import { useCallback, useState, type FormEvent } from "react";
import type { OperationConfig } from "../types";
import { parseFieldValue, validateOperation } from "../validation/validateOperation";
import { postCalculate } from "../api/calculatorApi";

type Status = "idle" | "loading" | "success" | "error";

export function useOperationForm(config: OperationConfig) {
  const initialValues = Object.fromEntries(
    config.fields.map((field) => [field.key, ""]),
  );

  const [values, setValues] = useState<Record<string, string>>(initialValues);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [result, setResult] = useState<number | null>(null);
  const [status, setStatus] = useState<Status>("idle");

  const handleChange = useCallback((key: string, raw: string) => {
    setValues((prev) => ({ ...prev, [key]: raw }));
    setFieldErrors((prev) => {
      if (!(key in prev)) return prev;
      const next = { ...prev };
      delete next[key];
      return next;
    });
    setFormError(null);
    setResult(null);
    setStatus("idle");
  }, []);

  const handleSubmit = useCallback(
    async (e: FormEvent) => {
      e.preventDefault();

      const parsed: Record<string, number> = {};
      const nextFieldErrors: Record<string, string> = {};

      for (const field of config.fields) {
        const value = parseFieldValue(values[field.key] ?? "");
        if (value === null) {
          nextFieldErrors[field.key] = "Enter a valid number.";
        } else {
          parsed[field.key] = value;
        }
      }

      if (Object.keys(nextFieldErrors).length > 0) {
        setFieldErrors(nextFieldErrors);
        setFormError(null);
        setResult(null);
        return;
      }

      const domainError = validateOperation(config.id, parsed);
      if (domainError) {
        setFieldErrors({});
        setFormError(domainError);
        setResult(null);
        return;
      }

      setFieldErrors({});
      setFormError(null);
      setStatus("loading");

      try {
        const value = await postCalculate(config.endpoint, parsed);
        setResult(value);
        setStatus("success");
      } catch (err) {
        setFormError(err instanceof Error ? err.message : "Calculation failed.");
        setStatus("error");
      }
    },
    [config, values],
  );

  return {
    values,
    fieldErrors,
    formError,
    result,
    status,
    handleChange,
    handleSubmit,
  };
}
