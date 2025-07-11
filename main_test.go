package main

import (
	"testing"
)

func TestEvalExprAdvanced(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
		err      bool
	}{
		{"2+2*3", 8, false},
		{"(2+2)*3", 12, false},
		{"-5+3", -2, false},
		{"-(2+3)*4", -20, false},
		{"2*(3+4*2)", 22, false},
		{"((2+3)*2)-4/2", 8, false},
		{"-(-5)", 5, false},
		{"2+", 0, true},
		{"*2+3", 0, true},
		{"2++3", 0, true},
		{"2+3)", 0, true},
		{"(2+3", 0, true},
		{"2+abc", 0, true},
		{"", 0, true},
		{"2/0", 2, false},         // división por cero: devuelve el número original
		{"(2+3)/(2-2)", 5, false}, // división por cero: devuelve el numerador
	}
	for _, test := range tests {
		res, err := evalExpr(test.expr)
		if test.err {
			if err == nil {
				t.Errorf("Expected error for expr '%s', got result %v", test.expr, res)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for expr '%s': %v", test.expr, err)
			} else if res != test.expected {
				t.Errorf("Expr '%s': expected %v, got %v", test.expr, test.expected, res)
			}
		}
	}
}
