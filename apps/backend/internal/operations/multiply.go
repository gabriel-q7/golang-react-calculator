package operations

type multiplyOperation struct{}

// NewMultiply returns the multiplication operation: a * b.
func NewMultiply() Operation { return multiplyOperation{} }

func (multiplyOperation) Name() string       { return "multiply" }
func (multiplyOperation) Operands() []string { return []string{"a", "b"} }
func (multiplyOperation) Apply(values map[string]float64) (float64, error) {
	return values["a"] * values["b"], nil
}
