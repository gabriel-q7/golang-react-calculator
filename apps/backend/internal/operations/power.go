package operations

import (
	"errors"
	"math"
)

// ErrUndefinedResult is returned when an operation has no real-valued
// result, such as a negative base raised to a fractional exponent.
var ErrUndefinedResult = errors.New("operation has no defined real result")

type powerOperation struct{}

// NewPower returns the exponentiation operation: base ^ exponent.
func NewPower() Operation { return powerOperation{} }

func (powerOperation) Name() string       { return "power" }
func (powerOperation) Operands() []string { return []string{"base", "exponent"} }
func (powerOperation) Apply(values map[string]float64) (float64, error) {
	base, exponent := values["base"], values["exponent"]
	if base == 0 && exponent < 0 {
		return 0, ErrDivisionByZero
	}

	result := math.Pow(base, exponent)
	if math.IsNaN(result) {
		return 0, ErrUndefinedResult
	}
	return result, nil
}
