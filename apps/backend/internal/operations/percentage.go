package operations

type percentageOperation struct{}

// NewPercentage returns the percentage operation: percent% of value.
func NewPercentage() Operation { return percentageOperation{} }

func (percentageOperation) Name() string       { return "percentage" }
func (percentageOperation) Operands() []string { return []string{"value", "percent"} }
func (percentageOperation) Apply(values map[string]float64) (float64, error) {
	return values["value"] * values["percent"] / 100, nil
}
