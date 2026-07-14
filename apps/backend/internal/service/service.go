// Package service orchestrates calculator operations for the API layer:
// it resolves an operation by name, validates that its operands are
// finite numbers, and delegates the actual math to the resolved
// operations.Operation. It has no knowledge of HTTP or JSON.
package service

import (
	"errors"
	"math"

	"calculator-backend/internal/operations"
)

// ErrInvalidInput is returned when an operand is not a finite number
// (e.g. NaN or +/-Infinity).
var ErrInvalidInput = errors.New("input must be a finite number")

// ErrUnknownOperation is returned when no operation is registered under
// the requested name.
var ErrUnknownOperation = errors.New("unknown operation")

// OperationResolver looks up an operation by name. CalculatorService
// depends on this small abstraction rather than *operations.Registry
// directly (dependency inversion): anything that can resolve a name to
// an operations.Operation works, including a test double.
type OperationResolver interface {
	Get(name string) (operations.Operation, bool)
}

// CalculatorService executes calculator operations by name.
type CalculatorService struct {
	resolver OperationResolver
}

// New creates a CalculatorService backed by resolver.
func New(resolver OperationResolver) *CalculatorService {
	return &CalculatorService{resolver: resolver}
}

// Execute resolves name to an operation and applies it to values. Every
// value is checked for finiteness first — a rule that applies uniformly
// to any operation, so it lives here rather than being repeated in each
// operations.Operation implementation. Adding a new operation never
// requires changing this method.
func (s *CalculatorService) Execute(name string, values map[string]float64) (float64, error) {
	op, ok := s.resolver.Get(name)
	if !ok {
		return 0, ErrUnknownOperation
	}

	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, ErrInvalidInput
		}
	}

	return op.Apply(values)
}
