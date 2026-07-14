package operations

import "fmt"

// Registry resolves operation names to Operation implementations. It is
// the single place operations get wired together — adding a new
// operation means implementing Operation and adding it to Default below;
// nothing that consumes a Registry (the service layer, the HTTP layer)
// needs to change.
type Registry struct {
	byName  map[string]Operation
	ordered []Operation
}

// NewRegistry builds a Registry from the given operations. It fails if
// two operations share a name, since that would make route/name
// resolution ambiguous.
func NewRegistry(ops ...Operation) (*Registry, error) {
	r := &Registry{byName: make(map[string]Operation, len(ops))}
	for _, op := range ops {
		if _, exists := r.byName[op.Name()]; exists {
			return nil, fmt.Errorf("operations: duplicate operation name %q", op.Name())
		}
		r.byName[op.Name()] = op
		r.ordered = append(r.ordered, op)
	}
	return r, nil
}

// Get resolves an operation by name.
func (r *Registry) Get(name string) (Operation, bool) {
	op, ok := r.byName[name]
	return op, ok
}

// All returns every registered operation, in registration order. Callers
// that mutate the result (e.g. sorting) get their own copy.
func (r *Registry) All() []Operation {
	return append([]Operation(nil), r.ordered...)
}

// Default returns the registry of every built-in calculator operation.
// This is the one function that has to know about all of them; adding an
// eighth operation is a one-line addition here plus its own new file —
// the service and HTTP layers pick it up automatically.
func Default() (*Registry, error) {
	return NewRegistry(
		NewAdd(),
		NewSubtract(),
		NewMultiply(),
		NewDivide(),
		NewPower(),
		NewSqrt(),
		NewPercentage(),
	)
}
