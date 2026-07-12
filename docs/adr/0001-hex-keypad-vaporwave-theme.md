# 0001. Honeycomb keypad and vaporwave theme for the calculator UI

## Status

Accepted

## Context

The frontend previously exposed the backend's seven operations
(add/subtract/multiply/divide/power/sqrt/percentage) as a tabbed form: one
tab per operation, each with labeled numeric inputs and a "Calculate"
button (`components/ui/tabs.tsx` + the now-removed `OperationCard.tsx` /
`useOperationForm.ts`). It worked, but it read as a settings form, not a
calculator — there was no running expression, no single keypad, and no
visual identity.

The brief asked for a redesign that (a) looks and behaves like a physical
calculator, with a display showing the running expression and result, (b)
replaces the per-operation form with a **single connected keypad** where
every key is a hexagon sharing edges with its neighbors (a honeycomb, not
a grid of separately-spaced buttons), and (c) applies a vaporwave visual
identity (dark background, neon pink/cyan/purple/blue, glow, animated
hover/press) without regressing backend integration, accessibility, or
test coverage.

The backend contract did not change and was not a degree of freedom here:
seven REST endpoints, each taking specific named fields (`a`/`b`,
`base`/`exponent`, `value`, `value`/`percent`) and returning `{result}` or
`{error}` (see [../api.md](../api.md)). Any UI shape had to keep mapping
user input onto exactly those calls.

## Decision

### Interaction model: physical calculator, not a form

Replace the tabbed form with a state machine (`useCalculatorEngine`, a
`useReducer` over `engine/reducer.ts`) that mimics a physical calculator:
digits build up a display buffer; pressing a binary operator (`+ − × ÷ ^
%`) locks in operand A and waits for operand B; `=` sends both operands to
the matching endpoint; `√` is unary and evaluates immediately against
whatever is currently displayed, without disturbing a pending binary
operation (so `2 × 16 √ =` computes `2 × √16 = 8`). This keeps the
existing per-operation endpoint/field mapping (`config.ts`,
`validateOperation.ts`, `calculatorApi.ts`) completely intact — reused
as-is from `engine/operations.ts` — so there is no backend or contract
change, only a different front-end shape wrapped around the same calls.

### Layout: a real honeycomb, not a button grid with rounded corners

Each key is a `<button>` clipped to a hexagon (`clip-path: polygon(50% 0%,
100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)`, "pointy-top" orientation).
21 keys are arranged as 7 rows of 3, with rows alternately offset by half
a hex-width and pulled together vertically (row step = 75% of hex height)
via CSS custom properties (`--row`, `--col`, `--stagger` set per button,
consumed in `HexKeypad.css`). That specific offset is what makes adjacent
hexagons interlock edge-to-edge into a honeycomb, rather than sitting in a
grid with gaps or with square hit-areas overlapping their neighbors'
visible corners. Modern browsers hit-test against the clipped shape, not
the bounding box, so the pointed corners of one key don't steal clicks
from the key next to it.

Digits, binary operators, "power-tier" operators (`^ % √`), controls
(`C ⌫ .`), and `=` each get a distinct neon accent (cyan / pink / purple /
blue respectively) via a `--key-color` custom property, which both
satisfies the "all four accents present" requirement and gives the keypad
a functional color code (you can tell operator classes apart at a glance).

### Visual identity: vaporwave, built as CSS on top of shadcn tokens

Tailwind v4 + shadcn/ui were already wired into the project
(`components.json`, `@theme inline` in `index.css`). Rather than
introducing a parallel design system, the existing shadcn color tokens
(`--background`, `--foreground`, `--card`, `--destructive`, etc.) were
repointed to a fixed dark vaporwave palette, and four new tokens
(`--neon-pink`, `--neon-cyan`, `--neon-purple`, `--neon-blue`) were
registered the same way the existing tokens are, so they're available
both as CSS variables and as Tailwind utilities (`text-neon-cyan`, etc).
Existing shadcn primitives (`Card`, `Alert`) pick up the new palette for
free with no changes to their source. The hexagon-specific rendering
(clip-path, glow, per-key color, hover/press/focus animation) lives in a
dedicated `HexKeypad.css`, imported only by `HexKeypad.tsx` — Tailwind
utilities handle everything that isn't hex-shaped.

Buttons use a **dark, semi-transparent fill with bright neon text/border**
rather than a solid neon fill, both because it reads as "glowing keys on a
dark panel" (the vaporwave look) and because neon-text-on-near-black
contrast is easy to keep at or above WCAG AA (~4.5:1) — see
"Accessibility" below. `prefers-reduced-motion: reduce` disables the
scale transforms on hover/press, keeping only the (already brief) color/
glow transition, so the animation requirement doesn't trap users who've
asked to avoid motion.

## Alternatives considered

1. **Keep the tabs/form UI, only re-skin the colors.**
   Rejected: the brief specifically asks for a single connected keypad and
   a running-expression display, which the one-operation-per-tab model
   can't represent — there's no notion of "the current calculation in
   progress" in a form.

2. **A conventional rectangular button grid (rounded squares) in a
   vaporwave palette, no hexagons.**
   Simpler and slightly less CSS (no `clip-path`, no offset-row math,
   no custom hit-testing concerns). Rejected because the brief explicitly
   requires *contiguous, edge-sharing hexagons*, and a rounded-square grid
   doesn't read as a honeycomb no matter the color treatment.

3. **SVG-based hex buttons (each key as an `<svg><polygon>`) instead of
   `clip-path` on a `<button>`.**
   SVG gives pixel-perfect control and arguably easier stroke/glow
   filters (`<feGaussianBlur>`), but turns every key into a non-native
   interactive element needing manual keyboard/focus wiring to behave
   like a button, and complicates using plain `<button>` semantics
   (disabled state, `:focus-visible`, form participation) that we get for
   free with `clip-path` on a real `<button>`. `clip-path` was chosen to
   keep each key a genuine, natively-accessible `<button>`.

4. **A CSS Grid / flexbox hex layout using only margin tricks (no
   absolute positioning), or a third-party hex-grid library.**
   Margin-collapse hex tricks are brittle across row-count and
   hex-size changes, and a library is a dependency for ~40 lines of
   geometry. Absolute positioning driven by simple per-button CSS custom
   properties (`--row`/`--col`/`--stagger`) was chosen as the smallest,
   most inspectable implementation, and it stays responsive via a
   `clamp()`-based `--hex-w` with no JS-computed layout.

5. **Adaptive theme (respect `prefers-color-scheme`, light + dark
   vaporwave variants).**
   Rejected: vaporwave's identity *is* the dark background with neon
   accents; a "light vaporwave" variant would need a materially different
   palette to keep contrast, doubling the design surface for a look whose
   whole point is the dark theme. The calculator is a self-contained
   themed app, not a general UI component library, so a fixed theme was
   judged the right scope.

6. **Solid neon fills on every key (pink/cyan/purple button backgrounds)
   instead of dark fill + neon border/text/glow.**
   More saturated and arguably more "vaporwave poster"-like at a glance,
   but multiple full-saturation neon fills next to each other are harder
   to keep AA-contrast-safe for their labels (light text on a mid-luminance
   saturated color is a much narrower safe range than light text on
   near-black), and a keypad of 21 solid-neon hexagons risks looking
   noisy rather than "glowing in the dark." Rejected in favor of the
   dark-fill-plus-glow treatment.

## Accessibility notes

- Every key is a real `<button>` (not a styled `<div>`), so it's reachable
  by keyboard (Tab) and activates on Enter/Space with no extra wiring.
- Symbol-only keys (`÷ × − ^ √ % ⌫ C .` `=`) carry an explicit `aria-label`
  (e.g. `aria-label="square root"`) so screen readers announce the
  operation, not just a glyph; the visible glyph itself is
  `aria-hidden="true"` to avoid double announcement.
- `clip-path` clips outlines along with everything else, so the default
  focus ring would be invisible on a hex button. Focus-visible instead
  uses a layered `filter: drop-shadow(...)` (white inner glow + colored
  outer glow), which follows the hex silhouette and gives a clearly
  visible, WCAG 2.4.7-compliant focus indicator without needing a
  shape-aware outline hack.
- The result line (`aria-live="polite"`) announces new values; the
  expression line above it is `aria-hidden` to avoid narrating every
  keystroke twice. Errors render through the existing shadcn `Alert`
  (`role="alert"`, assertive), unchanged from the old UI.
- Color pairings were hand-picked, not generated, specifically to clear
  ~4.5:1 (WCAG AA, normal text) against their backgrounds: white-lavender
  foreground (`#f4f1ff`) on the near-black page background (`#0b0618`) is
  >15:1; each neon key color against the dark button fill (`~#0b0618`–
  `#180f33`) is roughly 5:1–10:1 depending on hue (cyan brightest, pink/
  purple lowest but still ≥5:1); destructive text (`#1a0410`) on the
  destructive fill (`#ff5c8a`) is ~6.7:1.
- All hover/press scale animation respects `prefers-reduced-motion:
  reduce` (transform disabled, color/glow transition kept).

## Consequences

- The keypad's geometry (`HexKeypad.css`) is hand-tuned CSS custom
  properties rather than a reusable "hex grid" abstraction; adding an
  8th operation would mean re-deriving the row/stagger layout by hand,
  which is acceptable for a fixed 21-key calculator but would need
  revisiting if the keypad ever needed to be dynamically sized.
- `engine/operations.ts` + `engine/reducer.ts` + `hooks/useCalculatorEngine.ts`
  are new, directly-testable modules (pure reducer, pure operator mapping,
  a hook exercised via `renderHook`), replacing the deleted
  `useOperationForm.ts` — the old operation-config source of truth
  (`config.ts`, `validateOperation.ts`, `calculatorApi.ts`) was reused
  unchanged, so the backend contract has zero drift risk from this
  redesign.
- The old tabbed UI (`OperationCard.tsx`, its test, `useOperationForm.ts`)
  was deleted rather than kept alongside the new one, since the brief
  asks for a replacement, not an alternate mode.
