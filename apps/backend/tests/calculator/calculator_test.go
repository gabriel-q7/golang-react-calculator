package calculator_test

import (
	"errors"
	"math"
	"testing"

	"calculator-backend/internal/calculator"
)

func TestAdd(t *testing.T) {
	if got := calculator.Add(2, 3); got != 5 {
		t.Fatalf("Add() = %v, want 5", got)
	}
}

func TestSubtract(t *testing.T) {
	if got := calculator.Subtract(5, 3); got != 2 {
		t.Fatalf("Subtract() = %v, want 2", got)
	}
}

func TestMultiply(t *testing.T) {
	if got := calculator.Multiply(4, 3); got != 12 {
		t.Fatalf("Multiply() = %v, want 12", got)
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"basic", 9, 3, 3, nil},
		{"negative divisor", 9, -3, -3, nil},
		{"fractional result", 1, 4, 0.25, nil},
		{"by zero", 1, 0, 0, calculator.ErrDivisionByZero},
		{"zero by zero", 0, 0, 0, calculator.ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Divide(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Divide() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
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
		{"zero base negative exponent", 0, -1, 0, calculator.ErrDivisionByZero},
		{"negative base fractional exponent", -4, 0.5, 0, calculator.ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Power(tt.base, tt.exponent)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Power() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Power() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		want    float64
		wantErr error
	}{
		{"perfect square", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"non-perfect square", 2, math.Sqrt(2), nil},
		{"negative", -4, 0, calculator.ErrNegativeSqrt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Sqrt(tt.value)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Sqrt() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Sqrt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
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
			if got := calculator.Percentage(tt.value, tt.percent); got != tt.want {
				t.Fatalf("Percentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
