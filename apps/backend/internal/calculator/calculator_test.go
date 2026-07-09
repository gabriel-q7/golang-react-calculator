package calculator

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		op      Operator
		want    float64
		wantErr error
	}{
		{"add", 2, 3, Add, 5, nil},
		{"subtract", 5, 3, Subtract, 2, nil},
		{"multiply", 4, 3, Multiply, 12, nil},
		{"divide", 9, 3, Divide, 3, nil},
		{"divide by zero", 1, 0, Divide, 0, ErrDivisionByZero},
		{"unsupported operator", 1, 2, Operator("mod"), 0, ErrUnsupportedOperator{Operator: "mod"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.a, tt.b, tt.op)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
