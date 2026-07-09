// Package service orchestrates calculator operations for the API layer:
// it validates request-level input (finiteness) and delegates the actual
// math to the calculator package, translating domain errors as needed.
package service

import (
	"errors"
	"math"

	"calculator-backend/internal/calculator"
)

// ErrInvalidInput is returned when an operand is not a finite number
// (e.g. NaN or +/-Infinity).
var ErrInvalidInput = errors.New("input must be a finite number")

// CalculatorService exposes the calculator operations used by the API.
type CalculatorService struct{}

// New creates a CalculatorService.
func New() *CalculatorService {
	return &CalculatorService{}
}

func validate(values ...float64) error {
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return ErrInvalidInput
		}
	}
	return nil
}

// Add returns a + b.
func (s *CalculatorService) Add(a, b float64) (float64, error) {
	if err := validate(a, b); err != nil {
		return 0, err
	}
	return calculator.Add(a, b), nil
}

// Subtract returns a - b.
func (s *CalculatorService) Subtract(a, b float64) (float64, error) {
	if err := validate(a, b); err != nil {
		return 0, err
	}
	return calculator.Subtract(a, b), nil
}

// Multiply returns a * b.
func (s *CalculatorService) Multiply(a, b float64) (float64, error) {
	if err := validate(a, b); err != nil {
		return 0, err
	}
	return calculator.Multiply(a, b), nil
}

// Divide returns a / b.
func (s *CalculatorService) Divide(a, b float64) (float64, error) {
	if err := validate(a, b); err != nil {
		return 0, err
	}
	return calculator.Divide(a, b)
}

// Power returns base raised to the power of exponent.
func (s *CalculatorService) Power(base, exponent float64) (float64, error) {
	if err := validate(base, exponent); err != nil {
		return 0, err
	}
	return calculator.Power(base, exponent)
}

// Sqrt returns the square root of value.
func (s *CalculatorService) Sqrt(value float64) (float64, error) {
	if err := validate(value); err != nil {
		return 0, err
	}
	return calculator.Sqrt(value)
}

// Percentage returns percent% of value.
func (s *CalculatorService) Percentage(value, percent float64) (float64, error) {
	if err := validate(value, percent); err != nil {
		return 0, err
	}
	return calculator.Percentage(value, percent), nil
}
