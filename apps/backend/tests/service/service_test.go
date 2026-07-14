package service_test

import (
	"errors"
	"math"
	"testing"

	"calculator-backend/internal/operations"
	"calculator-backend/internal/service"
)

// fakeOperation lets the service tests below exercise CalculatorService's
// own logic (operand validation, error propagation, unknown-operation
// handling) without depending on any real operation's math — a stand-in
// that satisfies operations.Operation, decoupled from internal/operations
// entirely.
type fakeOperation struct {
	name    string
	applyFn func(values map[string]float64) (float64, error)
}

func (f fakeOperation) Name() string       { return f.name }
func (f fakeOperation) Operands() []string { return []string{"x"} }
func (f fakeOperation) Apply(values map[string]float64) (float64, error) {
	return f.applyFn(values)
}

type fakeResolver map[string]operations.Operation

func (r fakeResolver) Get(name string) (operations.Operation, bool) {
	op, ok := r[name]
	return op, ok
}

func TestCalculatorService_Execute_DelegatesToTheResolvedOperation(t *testing.T) {
	echo := fakeOperation{name: "echo", applyFn: func(v map[string]float64) (float64, error) {
		return v["x"] * 2, nil
	}}
	svc := service.New(fakeResolver{"echo": echo})

	got, err := svc.Execute("echo", map[string]float64{"x": 21})
	if err != nil || got != 42 {
		t.Fatalf("Execute() = %v, %v, want 42, nil", got, err)
	}
}

func TestCalculatorService_Execute_PropagatesOperationErrors(t *testing.T) {
	boom := errors.New("boom")
	failing := fakeOperation{name: "failing", applyFn: func(map[string]float64) (float64, error) {
		return 0, boom
	}}
	svc := service.New(fakeResolver{"failing": failing})

	_, err := svc.Execute("failing", map[string]float64{"x": 1})
	if !errors.Is(err, boom) {
		t.Fatalf("Execute() error = %v, want %v", err, boom)
	}
}

func TestCalculatorService_Execute_UnknownOperation(t *testing.T) {
	svc := service.New(fakeResolver{})

	_, err := svc.Execute("does-not-exist", nil)
	if !errors.Is(err, service.ErrUnknownOperation) {
		t.Fatalf("Execute() error = %v, want %v", err, service.ErrUnknownOperation)
	}
}

func TestCalculatorService_Execute_RejectsNonFiniteOperandsBeforeApplying(t *testing.T) {
	called := false
	op := fakeOperation{name: "op", applyFn: func(map[string]float64) (float64, error) {
		called = true
		return 0, nil
	}}
	svc := service.New(fakeResolver{"op": op})

	cases := map[string]float64{
		"NaN":  math.NaN(),
		"+Inf": math.Inf(1),
		"-Inf": math.Inf(-1),
	}
	for name, v := range cases {
		t.Run(name, func(t *testing.T) {
			called = false
			_, err := svc.Execute("op", map[string]float64{"x": v})
			if !errors.Is(err, service.ErrInvalidInput) {
				t.Fatalf("Execute() error = %v, want %v", err, service.ErrInvalidInput)
			}
			if called {
				t.Fatal("Execute() called Apply despite a non-finite operand")
			}
		})
	}
}

// TestCalculatorService_Execute_Integration exercises the service wired
// to the real operation registry, as it runs in production, covering
// domain errors that originate from a real operation rather than a fake.
func TestCalculatorService_Execute_Integration(t *testing.T) {
	registry, err := operations.Default()
	if err != nil {
		t.Fatalf("operations.Default() error = %v", err)
	}
	svc := service.New(registry)

	if got, err := svc.Execute("add", map[string]float64{"a": 2, "b": 3}); err != nil || got != 5 {
		t.Fatalf("Execute(add) = %v, %v, want 5, nil", got, err)
	}
	if got, err := svc.Execute("percentage", map[string]float64{"value": 200, "percent": 10}); err != nil || got != 20 {
		t.Fatalf("Execute(percentage) = %v, %v, want 20, nil", got, err)
	}
	if _, err := svc.Execute("divide", map[string]float64{"a": 1, "b": 0}); !errors.Is(err, operations.ErrDivisionByZero) {
		t.Fatalf("Execute(divide) error = %v, want %v", err, operations.ErrDivisionByZero)
	}
	if _, err := svc.Execute("sqrt", map[string]float64{"value": -4}); !errors.Is(err, operations.ErrNegativeSqrt) {
		t.Fatalf("Execute(sqrt) error = %v, want %v", err, operations.ErrNegativeSqrt)
	}
	if _, err := svc.Execute("power", map[string]float64{"base": -4, "exponent": 0.5}); !errors.Is(err, operations.ErrUndefinedResult) {
		t.Fatalf("Execute(power) error = %v, want %v", err, operations.ErrUndefinedResult)
	}
}
