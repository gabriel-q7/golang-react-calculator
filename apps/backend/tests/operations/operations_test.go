package operations_test

import (
	"errors"
	"math"
	"testing"

	"calculator-backend/internal/operations"
)

func TestAdd(t *testing.T) {
	op := operations.NewAdd()
	if op.Name() != "add" {
		t.Fatalf("Name() = %q, want %q", op.Name(), "add")
	}
	if got := op.Operands(); !equalStrings(got, []string{"a", "b"}) {
		t.Fatalf("Operands() = %v, want [a b]", got)
	}
	got, err := op.Apply(map[string]float64{"a": 2, "b": 3})
	if err != nil || got != 5 {
		t.Fatalf("Apply() = %v, %v, want 5, nil", got, err)
	}
}

func TestSubtract(t *testing.T) {
	op := operations.NewSubtract()
	got, err := op.Apply(map[string]float64{"a": 5, "b": 3})
	if err != nil || got != 2 {
		t.Fatalf("Apply() = %v, %v, want 2, nil", got, err)
	}
}

func TestMultiply(t *testing.T) {
	op := operations.NewMultiply()
	got, err := op.Apply(map[string]float64{"a": 4, "b": 3})
	if err != nil || got != 12 {
		t.Fatalf("Apply() = %v, %v, want 12, nil", got, err)
	}
}

func TestDivide(t *testing.T) {
	op := operations.NewDivide()
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"basic", 9, 3, 3, nil},
		{"negative divisor", 9, -3, -3, nil},
		{"fractional result", 1, 4, 0.25, nil},
		{"by zero", 1, 0, 0, operations.ErrDivisionByZero},
		{"zero by zero", 0, 0, 0, operations.ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Apply(map[string]float64{"a": tt.a, "b": tt.b})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Apply() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	op := operations.NewPower()
	tests := []struct {
		name           string
		base, exponent float64
		want           float64
		wantErr        error
	}{
		{"basic", 2, 10, 1024, nil},
		{"zero exponent", 5, 0, 1, nil},
		{"negative exponent", 2, -2, 0.25, nil},
		{"fractional exponent", 4, 0.5, 2, nil},
		{"zero base positive exponent", 0, 3, 0, nil},
		{"zero base negative exponent", 0, -1, 0, operations.ErrDivisionByZero},
		{"negative base fractional exponent", -4, 0.5, 0, operations.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Apply(map[string]float64{"base": tt.base, "exponent": tt.exponent})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Apply() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	op := operations.NewSqrt()
	tests := []struct {
		name    string
		value   float64
		want    float64
		wantErr error
	}{
		{"perfect square", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"non-perfect square", 2, math.Sqrt(2), nil},
		{"negative", -4, 0, operations.ErrNegativeSqrt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Apply(map[string]float64{"value": tt.value})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Apply() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	op := operations.NewPercentage()
	tests := []struct {
		name           string
		value, percent float64
		want           float64
	}{
		{"basic", 200, 10, 20},
		{"zero percent", 200, 0, 0},
		{"over 100 percent", 50, 150, 75},
		{"negative percent", 200, -10, -20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Apply(map[string]float64{"value": tt.value, "percent": tt.percent})
			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperandsContract checks every built-in operation reports the exact
// operand names the HTTP layer's JSON contract (docs/api.md) documents —
// a regression here would silently change the API's request shape.
func TestOperandsContract(t *testing.T) {
	tests := []struct {
		op       operations.Operation
		wantName string
		wantOps  []string
	}{
		{operations.NewAdd(), "add", []string{"a", "b"}},
		{operations.NewSubtract(), "subtract", []string{"a", "b"}},
		{operations.NewMultiply(), "multiply", []string{"a", "b"}},
		{operations.NewDivide(), "divide", []string{"a", "b"}},
		{operations.NewPower(), "power", []string{"base", "exponent"}},
		{operations.NewSqrt(), "sqrt", []string{"value"}},
		{operations.NewPercentage(), "percentage", []string{"value", "percent"}},
	}

	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			if got := tt.op.Name(); got != tt.wantName {
				t.Fatalf("Name() = %q, want %q", got, tt.wantName)
			}
			if got := tt.op.Operands(); !equalStrings(got, tt.wantOps) {
				t.Fatalf("Operands() = %v, want %v", got, tt.wantOps)
			}
		})
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
