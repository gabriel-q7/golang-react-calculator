package operations

import "errors"

// ErrDivisionByZero is returned when an operation would require dividing
// by zero — division itself, or exponentiation of zero to a negative
// exponent (see power.go).
var ErrDivisionByZero = errors.New("division by zero")

type divideOperation struct{}

// NewDivide returns the division operation: a / b.
func NewDivide() Operation { return divideOperation{} }

func (divideOperation) Name() string       { return "divide" }
func (divideOperation) Operands() []string { return []string{"a", "b"} }
func (divideOperation) Apply(values map[string]float64) (float64, error) {
	a, b := values["a"], values["b"]
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}
