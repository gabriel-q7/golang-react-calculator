// Package operations defines the calculator's operation abstraction and
// every concrete operation (add, subtract, multiply, divide, power, sqrt,
// percentage). Each operation is a small, independent type implementing
// Operation; nothing in this package or its callers needs to change when
// a new operation is added — see registry.go.
package operations

// Operation is a single calculator operation. It has no knowledge of
// HTTP, JSON, or how its inputs were obtained — it only knows its name,
// the named operands it needs, and how to compute a result from them.
//
// Every implementation must be safe to use as a value: Apply must not
// mutate shared state and must be safe for concurrent use, so any
// Operation is freely interchangeable with any other (Liskov
// substitution) wherever the Operation interface is expected.
type Operation interface {
	// Name identifies the operation (e.g. "add"). It's used as the
	// registry key and, by the HTTP layer, as the route "/api/<name>".
	Name() string

	// Operands lists the named inputs this operation expects, in the
	// order they should be documented/displayed. A binary operation like
	// Add returns []string{"a", "b"}; a unary one like Sqrt returns
	// []string{"value"}.
	Operands() []string

	// Apply computes the result from values, keyed by the names Operands
	// returns. Callers are expected to supply exactly those keys —
	// Apply itself only implements the mathematical domain rules (e.g.
	// division by zero), not request validation.
	Apply(values map[string]float64) (float64, error)
}
