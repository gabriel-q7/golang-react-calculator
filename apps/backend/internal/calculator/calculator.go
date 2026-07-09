// Package calculator implements the arithmetic core of the application.
package calculator

import "fmt"

// Operator identifies a supported arithmetic operation.
type Operator string

const (
	Add      Operator = "add"
	Subtract Operator = "subtract"
	Multiply Operator = "multiply"
	Divide   Operator = "divide"
)

// ErrDivisionByZero is returned when a Divide operation has a zero divisor.
var ErrDivisionByZero = fmt.Errorf("division by zero")

// ErrUnsupportedOperator is returned when Operator is not one of the known constants.
type ErrUnsupportedOperator struct {
	Operator Operator
}

func (e ErrUnsupportedOperator) Error() string {
	return fmt.Sprintf("unsupported operator: %q", e.Operator)
}

// Calculate applies op to a and b, returning the result.
func Calculate(a, b float64, op Operator) (float64, error) {
	switch op {
	case Add:
		return a + b, nil
	case Subtract:
		return a - b, nil
	case Multiply:
		return a * b, nil
	case Divide:
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrUnsupportedOperator{Operator: op}
	}
}
