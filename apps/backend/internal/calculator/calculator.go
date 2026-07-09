// Package calculator implements the pure arithmetic core of the
// application. Functions here have no knowledge of HTTP, JSON, or
// validation policy — they only encode the mathematical domain rules
// (e.g. division by zero, negative square roots).
package calculator

import (
	"errors"
	"math"
)

// ErrDivisionByZero is returned when a division or exponentiation would
// require dividing by zero (e.g. 0 raised to a negative exponent).
var ErrDivisionByZero = errors.New("division by zero")

// ErrNegativeSqrt is returned when Sqrt is called with a negative operand.
var ErrNegativeSqrt = errors.New("cannot take the square root of a negative number")

// ErrUndefinedResult is returned when an operation has no real-valued
// result, such as a negative base raised to a fractional exponent.
var ErrUndefinedResult = errors.New("operation has no defined real result")

// Add returns a + b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns a - b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns a * b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns a / b, or ErrDivisionByZero if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Power returns base raised to the power of exponent.
func Power(base, exponent float64) (float64, error) {
	if base == 0 && exponent < 0 {
		return 0, ErrDivisionByZero
	}

	result := math.Pow(base, exponent)
	if math.IsNaN(result) {
		return 0, ErrUndefinedResult
	}
	return result, nil
}

// Sqrt returns the square root of value, or ErrNegativeSqrt if value is
// negative.
func Sqrt(value float64) (float64, error) {
	if value < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(value), nil
}

// Percentage returns percent% of value, i.e. (value * percent) / 100.
func Percentage(value, percent float64) float64 {
	return value * percent / 100
}
