package lambda

import (
	"testing"
)

func TestParseVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected Var
	}{
		{"x", Var{Name: "x"}},
		{"y", Var{Name: "y"}},
		{"foo", Var{Name: "foo"}},
		{"var123", Var{Name: "var123"}},
		{"_test", Var{Name: "_test"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			v, ok := result.(Var)
			if !ok {
				t.Fatalf("Parse(%q) = %T, want Var", tt.input, result)
			}

			if v.Name != tt.expected.Name {
				t.Errorf("Parse(%q) = Var{Name: %q}, want Var{Name: %q}", tt.input, v.Name, tt.expected.Name)
			}
		})
	}
}

func TestParseAbstraction(t *testing.T) {
	tests := []struct {
		input       string
		expectedStr string
	}{
		{"λx.x", "λx.x"},
		{"\\x.x", "λx.x"},
		{"λx.λy.x", "λx.λy.x"},
		{"λf.λx.f x", "λf.λx.f x"},
		{"λx.(λy.y)", "λx.λy.y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if result.String() != tt.expectedStr {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, result.String(), tt.expectedStr)
			}
		})
	}
}

func TestParseApplication(t *testing.T) {
	tests := []struct {
		input       string
		expectedStr string
	}{
		{"f x", "f x"},
		{"f x y", "f x y"},
		{"(f x) y", "f x y"},
		{"f (x y)", "f (x y)"},
		{"(λx.x) y", "(λx.x) y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if result.String() != tt.expectedStr {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, result.String(), tt.expectedStr)
			}
		})
	}
}

func TestParseComplex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedStr string
	}{
		{
			name:        "Identity",
			input:       "λx.x",
			expectedStr: "λx.x",
		},
		{
			name:        "K combinator",
			input:       "λx.λy.x",
			expectedStr: "λx.λy.x",
		},
		{
			name:        "S combinator",
			input:       "λx.λy.λz.x z (y z)",
			expectedStr: "λx.λy.λz.x z (y z)",
		},
		{
			name:        "Church numeral 0",
			input:       "λf.λx.x",
			expectedStr: "λf.λx.x",
		},
		{
			name:        "Church numeral 1",
			input:       "λf.λx.f x",
			expectedStr: "λf.λx.f x",
		},
		{
			name:        "Church numeral 2",
			input:       "λf.λx.f (f x)",
			expectedStr: "λf.λx.f (f x)",
		},
		{
			name:        "PLUS",
			input:       "λm.λn.λf.λx.m f (n f x)",
			expectedStr: "λm.λn.λf.λx.m f (n f x)",
		},
		{
			name:        "Application with abstraction",
			input:       "(λx.x) (λy.y)",
			expectedStr: "(λx.x) (λy.y)",
		},
		{
			name:        "Nested applications",
			input:       "f x y z",
			expectedStr: "f x y z",
		},
		{
			name:        "Mixed",
			input:       "(λx.λy.x) a b",
			expectedStr: "(λx.λy.x) a b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if result.String() != tt.expectedStr {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, result.String(), tt.expectedStr)
			}
		})
	}
}

func TestParseAndReduce(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedStr string
	}{
		{
			name:        "Identity applied",
			input:       "(λx.x) y",
			expectedStr: "y",
		},
		{
			name:        "K combinator applied",
			input:       "(λx.λy.x) a b",
			expectedStr: "a",
		},
		{
			name:        "Nested reduction",
			input:       "(λx.x) ((λy.y) z)",
			expectedStr: "z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			// Reduce to normal form
			reduced, _ := Reduce(result, 100)

			if reduced.String() != tt.expectedStr {
				t.Errorf("Parse(%q) reduced = %q, want %q", tt.input, reduced.String(), tt.expectedStr)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Incomplete abstraction", "λx."},
		{"Missing dot", "λx y"},
		{"Unclosed paren", "(λx.x"},
		{"Extra closing paren", "λx.x)"},
		{"Invalid character at start", "123"},
		{"Extra characters", "x y z)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func TestParseWithWhitespace(t *testing.T) {
	tests := []struct {
		input       string
		expectedStr string
	}{
		{"  x  ", "x"},
		{"λx . x", "λx.x"},
		{"  λx.  λy.  x  ", "λx.λy.x"},
		{"f  x  y", "f x y"},
		{"  (  f  x  )  y  ", "f x y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if result.String() != tt.expectedStr {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, result.String(), tt.expectedStr)
			}
		})
	}
}

func TestParseChurchNumerals(t *testing.T) {
	tests := []struct {
		n     int
		input string
	}{
		{0, "λf.λx.x"},
		{1, "λf.λx.f x"},
		{2, "λf.λx.f (f x)"},
		{3, "λf.λx.f (f (f x))"},
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.n)), func(t *testing.T) {
			parsed, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			// Convert parsed to integer and verify
			parsedInt := ToInt(parsed)
			if parsedInt != tt.n {
				t.Errorf("ToInt(Parse(%q)) = %d, want %d", tt.input, parsedInt, tt.n)
			}

			// Also check against ChurchNumeral
			expected := ChurchNumeral(tt.n)
			if parsed.String() != expected.String() {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, parsed.String(), expected.String())
			}
		})
	}
}

func TestParseArithmetic(t *testing.T) {
	// Test PLUS 1 2 = 3
	plusInput := "(λm.λn.λf.λx.m f (n f x)) (λf.λx.f x) (λf.λx.f (f x))"
	result, err := Parse(plusInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reduced, _ := Reduce(result, 1000)
	got := ToInt(reduced)
	want := 3

	if got != want {
		t.Errorf("PLUS 1 2 = %d, want %d", got, want)
	}
}

func TestParseBooleans(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"TRUE", "λx.λy.x", true},
		{"FALSE", "λx.λy.y", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			result := ToBool(parsed)
			if result != tt.expected {
				t.Errorf("ToBool(Parse(%q)) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseBackslashNotation(t *testing.T) {
	// Test that backslash works as lambda
	tests := []struct {
		input       string
		expectedStr string
	}{
		{"\\x.x", "λx.x"},
		{"\\x.\\y.x", "λx.λy.x"},
		{"(\\x.x) y", "(λx.x) y"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if result.String() != tt.expectedStr {
				t.Errorf("Parse(%q).String() = %q, want %q", tt.input, result.String(), tt.expectedStr)
			}
		})
	}
}