package service_test

import (
	"errors"
	"math"
	"testing"

	"calculator-backend/internal/calculator"
	"calculator-backend/internal/service"
)

func TestCalculatorService_Arithmetic(t *testing.T) {
	svc := service.New()

	if got, err := svc.Add(2, 3); err != nil || got != 5 {
		t.Fatalf("Add() = %v, %v, want 5, nil", got, err)
	}
	if got, err := svc.Subtract(5, 3); err != nil || got != 2 {
		t.Fatalf("Subtract() = %v, %v, want 2, nil", got, err)
	}
	if got, err := svc.Multiply(4, 3); err != nil || got != 12 {
		t.Fatalf("Multiply() = %v, %v, want 12, nil", got, err)
	}
	if got, err := svc.Divide(9, 3); err != nil || got != 3 {
		t.Fatalf("Divide() = %v, %v, want 3, nil", got, err)
	}
	if got, err := svc.Power(2, 10); err != nil || got != 1024 {
		t.Fatalf("Power() = %v, %v, want 1024, nil", got, err)
	}
	if got, err := svc.Sqrt(9); err != nil || got != 3 {
		t.Fatalf("Sqrt() = %v, %v, want 3, nil", got, err)
	}
	if got, err := svc.Percentage(200, 10); err != nil || got != 20 {
		t.Fatalf("Percentage() = %v, %v, want 20, nil", got, err)
	}
}

func TestCalculatorService_DomainErrors(t *testing.T) {
	svc := service.New()

	if _, err := svc.Divide(1, 0); !errors.Is(err, calculator.ErrDivisionByZero) {
		t.Fatalf("Divide() error = %v, want %v", err, calculator.ErrDivisionByZero)
	}
	if _, err := svc.Sqrt(-4); !errors.Is(err, calculator.ErrNegativeSqrt) {
		t.Fatalf("Sqrt() error = %v, want %v", err, calculator.ErrNegativeSqrt)
	}
	if _, err := svc.Power(-4, 0.5); !errors.Is(err, calculator.ErrUndefinedResult) {
		t.Fatalf("Power() error = %v, want %v", err, calculator.ErrUndefinedResult)
	}
}

func TestCalculatorService_InvalidInput(t *testing.T) {
	svc := service.New()

	cases := []struct {
		name string
		call func() (float64, error)
	}{
		{"add NaN", func() (float64, error) { return svc.Add(math.NaN(), 1) }},
		{"subtract +Inf", func() (float64, error) { return svc.Subtract(math.Inf(1), 1) }},
		{"multiply -Inf", func() (float64, error) { return svc.Multiply(math.Inf(-1), 1) }},
		{"divide NaN", func() (float64, error) { return svc.Divide(1, math.NaN()) }},
		{"power Inf exponent", func() (float64, error) { return svc.Power(2, math.Inf(1)) }},
		{"sqrt NaN", func() (float64, error) { return svc.Sqrt(math.NaN()) }},
		{"percentage Inf", func() (float64, error) { return svc.Percentage(math.Inf(1), 10) }},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.call(); !errors.Is(err, service.ErrInvalidInput) {
				t.Fatalf("error = %v, want %v", err, service.ErrInvalidInput)
			}
		})
	}
}
