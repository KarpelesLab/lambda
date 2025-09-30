package lambda

import (
	"testing"
)

func TestParseDigitConstants(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"_0", 0},
		{"_1", 1},
		{"_2", 2},
		{"_5", 5},
		{"_10", 10},
		{"_42", 42},
		{"_100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}

			result := ToInt(expr)
			if result != tt.expected {
				t.Errorf("Parse(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseNamedConstants(t *testing.T) {
	tests := []struct {
		input    string
		expected Object
	}{
		{"_I", I},
		{"_K", K},
		{"_S", S},
		{"_B", B},
		{"_C", C},
		{"_W", W},
		{"_U", U},
		{"_TRUE", TRUE},
		{"_FALSE", FALSE},
		{"_T", T},
		{"_F", F},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}

			if expr.String() != tt.expected.String() {
				t.Errorf("Parse(%q) = %s, want %s", tt.input, expr.String(), tt.expected.String())
			}
		})
	}
}

func TestParseConstantsInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"_PLUS _2 _3", 5},
		{"_MULT _4 _5", 20},
		{"_SUCC _7", 8},
		{"_POW _2 _3", 8},
		{"_PLUS _1 (_MULT _2 _3)", 7},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}

			result, _ := Reduce(expr, 1000)
			got := ToInt(result)
			if got != tt.expected {
				t.Errorf("Parse(%q) reduced to %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseBooleanConstants(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"_TRUE", true},
		{"_FALSE", false},
		{"_T", true},
		{"_F", false},
		{"_AND _TRUE _TRUE", true},
		{"_AND _TRUE _FALSE", false},
		{"_OR _TRUE _FALSE", true},
		{"_OR _FALSE _FALSE", false},
		{"_NOT _TRUE", false},
		{"_NOT _FALSE", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}

			result, _ := Reduce(expr, 1000)
			got := ToBool(result)
			if got != tt.expected {
				t.Errorf("Parse(%q) reduced to %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}