package lambda

import (
	"testing"
)

func TestVarString(t *testing.T) {
	v := Var{Name: "x"}
	if v.String() != "x" {
		t.Errorf("Expected 'x', got '%s'", v.String())
	}
}

func TestAbstractionString(t *testing.T) {
	// λx.x
	abs := Abstraction{Param: "x", Body: Var{Name: "x"}}
	if abs.String() != "λx.x" {
		t.Errorf("Expected 'λx.x', got '%s'", abs.String())
	}

	// λx.λy.x
	nested := Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body:  Var{Name: "x"},
		},
	}
	if nested.String() != "λx.λy.x" {
		t.Errorf("Expected 'λx.λy.x', got '%s'", nested.String())
	}
}

func TestApplicationString(t *testing.T) {
	// x y
	app := Application{Func: Var{Name: "x"}, Arg: Var{Name: "y"}}
	if app.String() != "x y" {
		t.Errorf("Expected 'x y', got '%s'", app.String())
	}

	// (λx.x) y
	absApp := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	if absApp.String() != "(λx.x) y" {
		t.Errorf("Expected '(λx.x) y', got '%s'", absApp.String())
	}

	// x (y z)
	nestedApp := Application{
		Func: Var{Name: "x"},
		Arg:  Application{Func: Var{Name: "y"}, Arg: Var{Name: "z"}},
	}
	if nestedApp.String() != "x (y z)" {
		t.Errorf("Expected 'x (y z)', got '%s'", nestedApp.String())
	}
}

func TestFreeVars(t *testing.T) {
	// x has free variable x
	v := Var{Name: "x"}
	fv := v.FreeVars()
	if !fv["x"] || len(fv) != 1 {
		t.Errorf("Expected {x}, got %v", fv)
	}

	// λx.x has no free variables
	abs := Abstraction{Param: "x", Body: Var{Name: "x"}}
	fv = abs.FreeVars()
	if len(fv) != 0 {
		t.Errorf("Expected {}, got %v", fv)
	}

	// λx.y has free variable y
	abs2 := Abstraction{Param: "x", Body: Var{Name: "y"}}
	fv = abs2.FreeVars()
	if !fv["y"] || len(fv) != 1 {
		t.Errorf("Expected {y}, got %v", fv)
	}

	// x y has free variables x and y
	app := Application{Func: Var{Name: "x"}, Arg: Var{Name: "y"}}
	fv = app.FreeVars()
	if !fv["x"] || !fv["y"] || len(fv) != 2 {
		t.Errorf("Expected {x, y}, got %v", fv)
	}
}

func TestSubstitute(t *testing.T) {
	// x[x := y] = y
	v := Var{Name: "x"}
	result := v.Substitute("x", Var{Name: "y"})
	if result.String() != "y" {
		t.Errorf("Expected 'y', got '%s'", result.String())
	}

	// z[x := y] = z
	v2 := Var{Name: "z"}
	result = v2.Substitute("x", Var{Name: "y"})
	if result.String() != "z" {
		t.Errorf("Expected 'z', got '%s'", result.String())
	}

	// (λx.x)[x := y] = λx.x (no substitution in bound variable)
	abs := Abstraction{Param: "x", Body: Var{Name: "x"}}
	result = abs.Substitute("x", Var{Name: "y"})
	if result.String() != "λx.x" {
		t.Errorf("Expected 'λx.x', got '%s'", result.String())
	}

	// (λy.x)[x := z] = λy.z
	abs2 := Abstraction{Param: "y", Body: Var{Name: "x"}}
	result = abs2.Substitute("x", Var{Name: "z"})
	if result.String() != "λy.z" {
		t.Errorf("Expected 'λy.z', got '%s'", result.String())
	}

	// (x y)[x := z] = z y
	app := Application{Func: Var{Name: "x"}, Arg: Var{Name: "y"}}
	result = app.Substitute("x", Var{Name: "z"})
	if result.String() != "z y" {
		t.Errorf("Expected 'z y', got '%s'", result.String())
	}
}

func TestSubstituteCaptureAvoidance(t *testing.T) {
	// (λy.x)[x := y] should rename y to avoid capture
	abs := Abstraction{Param: "y", Body: Var{Name: "x"}}
	result := abs.Substitute("x", Var{Name: "y"})

	// The result should not have y free in the body after substitution
	// It should be something like λy0.y
	if result.String() == "λy.y" {
		t.Errorf("Variable capture occurred: got '%s'", result.String())
	}
}

func TestAlphaConvert(t *testing.T) {
	// λx.x renamed to λy.y
	abs := Abstraction{Param: "x", Body: Var{Name: "x"}}
	result := abs.AlphaConvert("x", "y")
	if result.String() != "λy.y" {
		t.Errorf("Expected 'λy.y', got '%s'", result.String())
	}

	// λx.λx.x renamed outer x to y: λy.λx.x
	nested := Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "x",
			Body:  Var{Name: "x"},
		},
	}
	result = nested.AlphaConvert("x", "y")
	if result.String() != "λy.λy.y" {
		t.Errorf("Expected 'λy.λy.y', got '%s'", result.String())
	}
}

func TestBetaReduce(t *testing.T) {
	// (λx.x) y → y
	app := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	result, reduced := app.BetaReduce()
	if !reduced {
		t.Errorf("Expected reduction to occur")
	}
	if result.String() != "y" {
		t.Errorf("Expected 'y', got '%s'", result.String())
	}

	// (λx.λy.x) z → λy.z
	app2 := Application{
		Func: Abstraction{
			Param: "x",
			Body: Abstraction{
				Param: "y",
				Body:  Var{Name: "x"},
			},
		},
		Arg: Var{Name: "z"},
	}
	result, reduced = app2.BetaReduce()
	if !reduced {
		t.Errorf("Expected reduction to occur")
	}
	if result.String() != "λy.z" {
		t.Errorf("Expected 'λy.z', got '%s'", result.String())
	}

	// x (no reduction possible)
	v := Var{Name: "x"}
	result, reduced = v.BetaReduce()
	if reduced {
		t.Errorf("Expected no reduction")
	}
	if result.String() != "x" {
		t.Errorf("Expected 'x', got '%s'", result.String())
	}
}

func TestEtaConvert(t *testing.T) {
	// λx.(f x) → f (when x is not free in f)
	abs := Abstraction{
		Param: "x",
		Body: Application{
			Func: Var{Name: "f"},
			Arg:  Var{Name: "x"},
		},
	}
	result, converted := abs.EtaConvert()
	if !converted {
		t.Errorf("Expected η-conversion to occur")
	}
	if result.String() != "f" {
		t.Errorf("Expected 'f', got '%s'", result.String())
	}

	// λx.(x x) should not η-convert (x is free in the function part)
	abs2 := Abstraction{
		Param: "x",
		Body: Application{
			Func: Var{Name: "x"},
			Arg:  Var{Name: "x"},
		},
	}
	result, converted = abs2.EtaConvert()
	if converted {
		t.Errorf("Expected no η-conversion for λx.(x x)")
	}
}

func TestChurchNumeral(t *testing.T) {
	// 0 := λf.λx.x
	zero := ChurchNumeral(0)
	if zero.String() != "λf.λx.x" {
		t.Errorf("Expected 'λf.λx.x', got '%s'", zero.String())
	}

	// 1 := λf.λx.f x
	one := ChurchNumeral(1)
	if one.String() != "λf.λx.f x" {
		t.Errorf("Expected 'λf.λx.f x', got '%s'", one.String())
	}

	// 2 := λf.λx.f (f x)
	two := ChurchNumeral(2)
	if two.String() != "λf.λx.f (f x)" {
		t.Errorf("Expected 'λf.λx.f (f x)', got '%s'", two.String())
	}

	// 3 := λf.λx.f (f (f x))
	three := ChurchNumeral(3)
	if three.String() != "λf.λx.f (f (f x))" {
		t.Errorf("Expected 'λf.λx.f (f (f x))', got '%s'", three.String())
	}
}

func TestChurchNumeralBetaReduction(t *testing.T) {
	// Apply Church numeral 2 to a function and argument
	// 2 f x = (λf.λx.f (f x)) f x → (λx.f (f x)) x → f (f x)

	two := ChurchNumeral(2)

	// Apply to f
	step1 := Application{Func: two, Arg: Var{Name: "g"}}
	result, reduced := step1.BetaReduce()
	if !reduced {
		t.Errorf("Expected first β-reduction to occur")
	}

	// Apply to x
	step2 := Application{Func: result, Arg: Var{Name: "a"}}
	result, reduced = step2.BetaReduce()
	if !reduced {
		t.Errorf("Expected second β-reduction to occur")
	}

	// Should be g (g a)
	if result.String() != "g (g a)" {
		t.Errorf("Expected 'g (g a)', got '%s'", result.String())
	}
}

func TestComplexBetaReduction(t *testing.T) {
	// ((λx.λy.y) a) b → (λy.y) b → b

	term := Application{
		Func: Application{
			Func: Abstraction{
				Param: "x",
				Body: Abstraction{
					Param: "y",
					Body:  Var{Name: "y"},
				},
			},
			Arg: Var{Name: "a"},
		},
		Arg: Var{Name: "b"},
	}

	// First reduction: ((λx.λy.y) a) b → (λy.y) b
	result, reduced := term.BetaReduce()
	if !reduced {
		t.Errorf("Expected first reduction")
	}

	// Second reduction: (λy.y) b → b
	result, reduced = result.BetaReduce()
	if !reduced {
		t.Errorf("Expected second reduction")
	}

	if result.String() != "b" {
		t.Errorf("Expected 'b', got '%s'", result.String())
	}
}

func TestReduce(t *testing.T) {
	// Test the Reduce helper function
	// (λx.x) y should reduce to y in one step
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}

	result, steps := Reduce(term, 10)

	if steps != 1 {
		t.Errorf("Expected 1 reduction step, got %d", steps)
	}

	if result.String() != "y" {
		t.Errorf("Expected 'y', got '%s'", result.String())
	}

	// Test with default limit (0 means use default 1000)
	result, steps = Reduce(term, 0)
	if steps != 1 {
		t.Errorf("Expected 1 reduction step with default limit, got %d", steps)
	}

	// Test with no reduction possible
	v := Var{Name: "x"}
	result, steps = Reduce(v, 10)
	if steps != 0 {
		t.Errorf("Expected 0 reduction steps for variable, got %d", steps)
	}
}

func TestFreshVar(t *testing.T) {
	avoid := map[string]bool{"x": true, "x0": true, "x1": true}
	fresh := freshVar("x", avoid)
	if fresh != "x2" {
		t.Errorf("Expected 'x2', got '%s'", fresh)
	}

	avoid2 := map[string]bool{}
	fresh2 := freshVar("y", avoid2)
	if fresh2 != "y" {
		t.Errorf("Expected 'y', got '%s'", fresh2)
	}
}

func TestToInt(t *testing.T) {
	// Test basic Church numerals
	if ToInt(ChurchNumeral(0)) != 0 {
		t.Errorf("Expected 0, got %d", ToInt(ChurchNumeral(0)))
	}

	if ToInt(ChurchNumeral(1)) != 1 {
		t.Errorf("Expected 1, got %d", ToInt(ChurchNumeral(1)))
	}

	if ToInt(ChurchNumeral(5)) != 5 {
		t.Errorf("Expected 5, got %d", ToInt(ChurchNumeral(5)))
	}
}

func TestMult(t *testing.T) {
	// 2 * 3 = 6
	two := ChurchNumeral(2)
	three := ChurchNumeral(3)

	var result Term = Application{
		Func: Application{
			Func: MULT,
			Arg:  two,
		},
		Arg: three,
	}

	// Reduce the result
	result, _ = Reduce(result, 100)

	resultInt := ToInt(result)
	if resultInt != 6 {
		t.Errorf("Expected 2 * 3 = 6, got %d", resultInt)
	}
}

func TestPow(t *testing.T) {
	// 2^3 = 8
	two := ChurchNumeral(2)
	three := ChurchNumeral(3)

	var result Term = Application{
		Func: Application{
			Func: POW,
			Arg:  two,
		},
		Arg: three,
	}

	// Reduce the result
	result, _ = Reduce(result, 100)

	resultInt := ToInt(result)
	if resultInt != 8 {
		t.Errorf("Expected 2^3 = 8, got %d", resultInt)
	}
}

func TestFactorial(t *testing.T) {
	// Test factorial(3) = 6
	three := ChurchNumeral(3)

	// Apply FACTORIAL to 3
	var result Term = Application{
		Func: FACTORIAL,
		Arg:  three,
	}

	// We need to reduce this, but factorial involves the Y combinator
	// which can create infinite expansion. We'll limit reductions.
	// For factorial(3), we need enough reductions to compute the result.
	result, _ = Reduce(result, 1000)

	resultInt := ToInt(result)
	if resultInt != 6 {
		t.Errorf("Expected factorial(3) = 6, got %d", resultInt)
		t.Logf("Result term: %s", result.String())
	}
}

func TestFAC(t *testing.T) {
	// Test FAC 3 = 6 using alternative factorial implementation
	three := ChurchNumeral(3)

	// Apply FAC to 3
	var result Term = Application{
		Func: FAC,
		Arg:  three,
	}

	// Reduce
	result, _ = Reduce(result, 1000)

	resultInt := ToInt(result)
	if resultInt != 6 {
		t.Errorf("Expected FAC(3) = 6, got %d", resultInt)
		t.Logf("Result term: %s", result.String())
	}
}

func TestFIB(t *testing.T) {
	// Test Fibonacci sequence: 0, 1, 1, 2, 3, 5, 8, ...
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
	}

	for _, tt := range tests {
		n := ChurchNumeral(tt.n)

		// FIB n returns a Church numeral directly
		var result Term = Application{
			Func: FIB,
			Arg:  n,
		}

		// Reduce
		result, _ = Reduce(result, 2000)

		resultInt := ToInt(result)
		if resultInt != tt.expected {
			t.Errorf("Expected FIB(%d) = %d, got %d", tt.n, tt.expected, resultInt)
			t.Logf("Result term: %s", result.String())
		}
	}
}

func TestToBool(t *testing.T) {
	// Test TRUE
	if !ToBool(TRUE) {
		t.Error("Expected TRUE to convert to true")
	}

	// Test FALSE
	if ToBool(FALSE) {
		t.Error("Expected FALSE to convert to false")
	}
}

func TestAND(t *testing.T) {
	tests := []struct {
		a        Term
		b        Term
		expected bool
		name     string
	}{
		{TRUE, TRUE, true, "TRUE AND TRUE"},
		{TRUE, FALSE, false, "TRUE AND FALSE"},
		{FALSE, TRUE, false, "FALSE AND TRUE"},
		{FALSE, FALSE, false, "FALSE AND FALSE"},
	}

	for _, tt := range tests {
		var result Term = Application{
			Func: Application{
				Func: AND,
				Arg:  tt.a,
			},
			Arg: tt.b,
		}

		result, _ = Reduce(result, 100)
		resultBool := ToBool(result)

		if resultBool != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, resultBool)
		}
	}
}

func TestOR(t *testing.T) {
	tests := []struct {
		a        Term
		b        Term
		expected bool
		name     string
	}{
		{TRUE, TRUE, true, "TRUE OR TRUE"},
		{TRUE, FALSE, true, "TRUE OR FALSE"},
		{FALSE, TRUE, true, "FALSE OR TRUE"},
		{FALSE, FALSE, false, "FALSE OR FALSE"},
	}

	for _, tt := range tests {
		var result Term = Application{
			Func: Application{
				Func: OR,
				Arg:  tt.a,
			},
			Arg: tt.b,
		}

		result, _ = Reduce(result, 100)
		resultBool := ToBool(result)

		if resultBool != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, resultBool)
		}
	}
}

func TestNOT(t *testing.T) {
	// NOT TRUE = FALSE
	var result Term = Application{
		Func: NOT,
		Arg:  TRUE,
	}
	result, _ = Reduce(result, 100)
	if ToBool(result) {
		t.Error("Expected NOT TRUE to be false")
	}

	// NOT FALSE = TRUE
	result = Application{
		Func: NOT,
		Arg:  FALSE,
	}
	result, _ = Reduce(result, 100)
	if !ToBool(result) {
		t.Error("Expected NOT FALSE to be true")
	}
}