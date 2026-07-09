import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import type { OperationConfig } from "../types";
import { useOperationForm } from "../hooks/useOperationForm";

export interface OperationCardProps {
  config: OperationConfig;
}

export function OperationCard({ config }: OperationCardProps) {
  const { values, fieldErrors, formError, result, status, handleChange, handleSubmit } =
    useOperationForm(config);

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          {config.label}{" "}
          <span className="text-muted-foreground font-normal">
            ({config.symbol})
          </span>
        </CardTitle>
        <CardDescription>
          Enter {config.fields.length === 1 ? "a value" : "values"} and
          calculate {config.label.toLowerCase()}.
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="flex flex-col gap-4 sm:flex-row sm:flex-wrap">
          {config.fields.map((field) => {
            const inputId = `${config.id}-${field.key}`;
            const fieldError = fieldErrors[field.key];
            return (
              <div key={field.key} className="flex flex-1 flex-col gap-1.5 min-w-[8rem]">
                <Label htmlFor={inputId}>{field.label}</Label>
                <Input
                  id={inputId}
                  inputMode="decimal"
                  aria-label={field.label}
                  aria-invalid={Boolean(fieldError)}
                  value={values[field.key] ?? ""}
                  onChange={(e) => handleChange(field.key, e.target.value)}
                />
                {fieldError && (
                  <p className="text-xs text-destructive">{fieldError}</p>
                )}
              </div>
            );
          })}
        </CardContent>
        <CardFooter className="flex flex-col items-stretch gap-3">
          <Button type="submit" disabled={status === "loading"}>
            {status === "loading" ? "Calculating…" : "Calculate"}
          </Button>

          {formError && (
            <Alert variant="destructive">
              <AlertDescription>{formError}</AlertDescription>
            </Alert>
          )}

          {result !== null && (
            <p aria-live="polite" className="text-lg font-semibold">
              = {result}
            </p>
          )}
        </CardFooter>
      </form>
    </Card>
  );
}
