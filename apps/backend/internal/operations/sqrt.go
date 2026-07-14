package operations

import (
	"errors"
	"math"
)

// ErrNegativeSqrt is returned when Sqrt is applied to a negative operand.
var ErrNegativeSqrt = errors.New("cannot take the square root of a negative number")

type sqrtOperation struct{}

// NewSqrt returns the square root operation: √value.
func NewSqrt() Operation { return sqrtOperation{} }

func (sqrtOperation) Name() string       { return "sqrt" }
func (sqrtOperation) Operands() []string { return []string{"value"} }
func (sqrtOperation) Apply(values map[string]float64) (float64, error) {
	value := values["value"]
	if value < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(value), nil
}
