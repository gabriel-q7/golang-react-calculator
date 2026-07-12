import type { CSSProperties } from "react";
import { cn } from "@/lib/utils";

export type HexKeyVariant = "digit" | "operator" | "operator-alt" | "control" | "equals";

export interface HexButtonProps {
  label: string;
  ariaLabel: string;
  variant: HexKeyVariant;
  disabled?: boolean;
  onClick: () => void;
  /** Grid placement, forwarded as CSS custom properties by HexKeypad. */
  gridStyle: CSSProperties;
}

export function HexButton({
  label,
  ariaLabel,
  variant,
  disabled,
  onClick,
  gridStyle,
}: HexButtonProps) {
  return (
    <button
      type="button"
      className={cn("hex-btn", `hex-btn--${variant}`)}
      style={gridStyle}
      aria-label={ariaLabel}
      disabled={disabled}
      onClick={onClick}
    >
      <span aria-hidden="true">{label}</span>
    </button>
  );
}
