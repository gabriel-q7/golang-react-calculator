package operations_test

import (
	"testing"

	"calculator-backend/internal/operations"
)

type stubOperation struct{ name string }

func (s stubOperation) Name() string       { return s.name }
func (s stubOperation) Operands() []string { return []string{"x"} }
func (s stubOperation) Apply(values map[string]float64) (float64, error) {
	return values["x"], nil
}

func TestRegistry_GetAndAll(t *testing.T) {
	a, b := stubOperation{"a"}, stubOperation{"b"}
	reg, err := operations.NewRegistry(a, b)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	got, ok := reg.Get("a")
	if !ok || got.Name() != "a" {
		t.Fatalf("Get(%q) = %v, %v", "a", got, ok)
	}

	if _, ok := reg.Get("missing"); ok {
		t.Fatalf("Get(%q) found an operation that was never registered", "missing")
	}

	all := reg.All()
	if len(all) != 2 {
		t.Fatalf("All() returned %d operations, want 2", len(all))
	}
}

func TestRegistry_RejectsDuplicateNames(t *testing.T) {
	_, err := operations.NewRegistry(stubOperation{"dup"}, stubOperation{"dup"})
	if err == nil {
		t.Fatal("NewRegistry() with duplicate names: want error, got nil")
	}
}

func TestRegistry_AllReturnsACopy(t *testing.T) {
	reg, err := operations.NewRegistry(stubOperation{"a"})
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	all := reg.All()
	all[0] = stubOperation{"mutated"}

	if got, _ := reg.Get("a"); got.Name() != "a" {
		t.Fatalf("mutating All()'s result affected the registry: Get(%q).Name() = %q", "a", got.Name())
	}
}

// TestDefault_RegistersAllSevenOperations guards against a new operation
// being added under internal/operations but forgotten in Default() (the
// one place that has to know about every built-in operation).
func TestDefault_RegistersAllSevenOperations(t *testing.T) {
	reg, err := operations.Default()
	if err != nil {
		t.Fatalf("Default() error = %v", err)
	}

	want := []string{"add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"}
	all := reg.All()
	if len(all) != len(want) {
		t.Fatalf("Default() registered %d operations, want %d", len(all), len(want))
	}
	for _, name := range want {
		if _, ok := reg.Get(name); !ok {
			t.Fatalf("Default() registry is missing operation %q", name)
		}
	}
}
