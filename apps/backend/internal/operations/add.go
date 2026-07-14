package operations

type addOperation struct{}

// NewAdd returns the addition operation: a + b.
func NewAdd() Operation { return addOperation{} }

func (addOperation) Name() string       { return "add" }
func (addOperation) Operands() []string { return []string{"a", "b"} }
func (addOperation) Apply(values map[string]float64) (float64, error) {
	return values["a"] + values["b"], nil
}
