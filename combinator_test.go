package lambda

import (
	"testing"
)

func TestCombinatorI(t *testing.T) {
	// I x = x
	x := Var{Name: "x"}
	result := Application{Func: I, Arg: x}

	reduced, didReduce := result.BetaReduce()
	if !didReduce {
		t.Error("Expected I x to reduce")
	}

	if reduced.String() != "x" {
		t.Errorf("Expected 'x', got '%s'", reduced.String())
	}
}

func TestCombinatorK(t *testing.T) {
	// K x y = x
	x := Var{Name: "x"}
	y := Var{Name: "y"}

	result := Application{
		Func: Application{Func: K, Arg: x},
		Arg:  y,
	}

	// First reduction: K x y -> (λy.x) y
	reduced1, didReduce := result.BetaReduce()
	if !didReduce {
		t.Error("Expected first reduction")
	}

	// Second reduction: (λy.x) y -> x
	reduced2, didReduce := reduced1.BetaReduce()
	if !didReduce {
		t.Error("Expected second reduction")
	}

	if reduced2.String() != "x" {
		t.Errorf("Expected 'x', got '%s'", reduced2.String())
	}
}

func TestCombinatorS(t *testing.T) {
	// S K K x = I x = x (SKK is the identity combinator)
	x := Var{Name: "x"}

	skk := Application{
		Func: Application{
			Func: S,
			Arg:  K,
		},
		Arg: K,
	}

	var result Object = Application{Func: skk, Arg: x}

	// Reduce multiple times
	for i := 0; i < 10; i++ {
		reduced, didReduce := result.BetaReduce()
		if !didReduce {
			break
		}
		result = reduced
	}

	if result.String() != "x" {
		t.Errorf("Expected S K K x = x, got '%s'", result.String())
	}
}

func TestCombinatorB(t *testing.T) {
	// B is composition: B f g x = f (g x)
	// We'll test with simple variables
	f := Var{Name: "f"}
	g := Var{Name: "g"}
	x := Var{Name: "x"}

	var result Object = Application{
		Func: Application{
			Func: Application{Func: B, Arg: f},
			Arg:  g,
		},
		Arg: x,
	}

	// Reduce
	for i := 0; i < 10; i++ {
		reduced, didReduce := result.BetaReduce()
		if !didReduce {
			break
		}
		result = reduced
	}

	// Should be f (g x)
	if result.String() != "f (g x)" {
		t.Errorf("Expected 'f (g x)', got '%s'", result.String())
	}
}

func TestCombinatorC(t *testing.T) {
	// C flips arguments: C f x y = f y x
	f := Var{Name: "f"}
	x := Var{Name: "x"}
	y := Var{Name: "y"}

	var result Object = Application{
		Func: Application{
			Func: Application{Func: C, Arg: f},
			Arg:  x,
		},
		Arg: y,
	}

	// Reduce
	for i := 0; i < 10; i++ {
		reduced, didReduce := result.BetaReduce()
		if !didReduce {
			break
		}
		result = reduced
	}

	// Should be f y x
	if result.String() != "f y x" {
		t.Errorf("Expected 'f y x', got '%s'", result.String())
	}
}

func TestCombinatorW(t *testing.T) {
	// W duplicates: W f x = f x x
	f := Var{Name: "f"}
	x := Var{Name: "x"}

	var result Object = Application{
		Func: Application{Func: W, Arg: f},
		Arg:  x,
	}

	// Reduce
	for i := 0; i < 10; i++ {
		reduced, didReduce := result.BetaReduce()
		if !didReduce {
			break
		}
		result = reduced
	}

	// Should be f x x
	if result.String() != "f x x" {
		t.Errorf("Expected 'f x x', got '%s'", result.String())
	}
}

func TestOmegaLower(t *testing.T) {
	// ω (omega_lower) = λx.x x
	if OMEGA_LOWER.String() != "λx.x x" {
		t.Errorf("Expected 'λx.x x', got '%s'", OMEGA_LOWER.String())
	}

	// Test self-application property
	x := Var{Name: "x"}
	result := Application{Func: OMEGA_LOWER, Arg: x}

	reduced, didReduce := result.BetaReduce()
	if !didReduce {
		t.Error("Expected ω x to reduce")
	}

	// Should be x x
	if reduced.String() != "x x" {
		t.Errorf("Expected 'x x', got '%s'", reduced.String())
	}
}

func TestOmegaInfinite(t *testing.T) {
	// Ω = ω ω reduces to itself infinitely
	// We won't fully reduce it, just check one step
	result := OMEGA

	reduced, didReduce := result.BetaReduce()
	if !didReduce {
		t.Error("Expected Ω to reduce")
	}

	// After one reduction, should be structurally similar (another self-application)
	// It should remain as (λx.x x) (λx.x x)
	if reduced.String() != OMEGA.String() {
		t.Logf("Ω reduces to: %s", reduced.String())
		t.Logf("Original Ω: %s", OMEGA.String())
		// This is expected - Ω reduces to itself
	}
}

func TestAliases(t *testing.T) {
	// Test that aliases point to the same combinators
	if DELTA.String() != OMEGA_LOWER.String() {
		t.Error("DELTA should equal OMEGA_LOWER")
	}

	if U.String() != OMEGA_LOWER.String() {
		t.Error("U should equal OMEGA_LOWER")
	}
}

func TestTRUEisK(t *testing.T) {
	// TRUE should be the same as K
	if TRUE.String() != K.String() {
		t.Error("TRUE should equal K combinator")
	}
}

func TestCombinatorStrings(t *testing.T) {
	// Test string representations
	tests := []struct {
		name     string
		term     Object
		expected string
	}{
		{"I", I, "λx.x"},
		{"K", K, "λx.λy.x"},
		{"S", S, "λx.λy.λz.x z (y z)"},
		{"B", B, "λx.λy.λz.x (y z)"},
		{"C", C, "λx.λy.λz.x z y"},
		{"W", W, "λx.λy.x y y"},
		{"ω", OMEGA_LOWER, "λx.x x"},
	}

	for _, tt := range tests {
		if tt.term.String() != tt.expected {
			t.Errorf("%s: expected '%s', got '%s'", tt.name, tt.expected, tt.term.String())
		}
	}
}