package operations

type subtractOperation struct{}

// NewSubtract returns the subtraction operation: a - b.
func NewSubtract() Operation { return subtractOperation{} }

func (subtractOperation) Name() string       { return "subtract" }
func (subtractOperation) Operands() []string { return []string{"a", "b"} }
func (subtractOperation) Apply(values map[string]float64) (float64, error) {
	return values["a"] - values["b"], nil
}
