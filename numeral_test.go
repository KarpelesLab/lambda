package lambda

import (
	"testing"
)

func TestNumeralString(t *testing.T) {
	if s := Numeral(0).String(); s != "[0]" {
		t.Errorf("Numeral(0).String() = %q, want [0]", s)
	}
	if s := Numeral(42).String(); s != "[42]" {
		t.Errorf("Numeral(42).String() = %q, want [42]", s)
	}
}

func TestNumeralFreeVars(t *testing.T) {
	fv := Numeral(5).FreeVars()
	if len(fv) != 0 {
		t.Errorf("Numeral should have no free vars, got %v", fv)
	}
}

func TestNumeralToInt(t *testing.T) {
	for _, n := range []int{0, 1, 2, 5, 10, 100} {
		got := ToInt(Numeral(uint64(n)))
		if got != n {
			t.Errorf("ToInt(Numeral(%d)) = %d", n, got)
		}
	}
}

func TestNumeralExpand(t *testing.T) {
	// Numeral(3).Expand() should equal ChurchNumeral(3)
	expanded := Numeral(3).Expand()
	church := ChurchNumeral(3)
	if expanded.String() != church.String() {
		t.Errorf("Numeral(3).Expand() = %s, want %s", expanded, church)
	}
}

func TestNumeralBetaReduce(t *testing.T) {
	// Numeral itself doesn't reduce
	_, reduced := Numeral(5).BetaReduce()
	if reduced {
		t.Error("Numeral should not beta-reduce on its own")
	}
}

func TestNumeralApplyPartial(t *testing.T) {
	// Numeral(3) f → NumeralApply{3, f}
	app := Application{Func: Numeral(3), Arg: Var{Name: "f"}}
	result, reduced := app.BetaReduce()
	if !reduced {
		t.Fatal("Numeral(3) f should reduce")
	}
	na, ok := result.(NumeralApply)
	if !ok {
		t.Fatalf("expected NumeralApply, got %T: %s", result, result)
	}
	if na.N != 3 {
		t.Errorf("NumeralApply.N = %d, want 3", na.N)
	}
}

func TestNumeralApplyFull(t *testing.T) {
	// Numeral(3) f x → f(f(f(x))) in two steps
	term := Application{
		Func: Application{Func: Numeral(3), Arg: Var{Name: "f"}},
		Arg:  Var{Name: "x"},
	}
	// Step 1: Numeral(3) f → NumeralApply
	result, reduced := term.BetaReduce()
	if !reduced {
		t.Fatal("step 1 should reduce")
	}
	// Step 2: NumeralApply{3, f} x → f(f(f(x)))
	result, reduced = result.BetaReduce()
	if !reduced {
		t.Fatal("step 2 should reduce")
	}
	// Verify structure: f(f(f(x)))
	expected := Application{
		Func: Var{Name: "f"},
		Arg: Application{
			Func: Var{Name: "f"},
			Arg: Application{
				Func: Var{Name: "f"},
				Arg:  Var{Name: "x"},
			},
		},
	}
	if result.String() != expected.String() {
		t.Errorf("got %s, want %s", result, expected)
	}
}

func TestNumeralZero(t *testing.T) {
	// Numeral(0) f x → x (zero applications of f)
	term := Application{
		Func: Application{Func: Numeral(0), Arg: Var{Name: "f"}},
		Arg:  Var{Name: "x"},
	}
	result, _ := Reduce(term, 10)
	if v, ok := result.(Var); !ok || v.Name != "x" {
		t.Errorf("Numeral(0) f x = %s, want x", result)
	}
}

func TestNumeralArithmetic(t *testing.T) {
	tests := []struct {
		name string
		term Term
		want int
	}{
		{
			"SUCC 3",
			Application{Func: SUCC, Arg: Numeral(3)},
			4,
		},
		{
			"PLUS 3 5",
			Application{Func: Application{Func: PLUS, Arg: Numeral(3)}, Arg: Numeral(5)},
			8,
		},
		{
			"MULT 3 4",
			Application{Func: Application{Func: MULT, Arg: Numeral(3)}, Arg: Numeral(4)},
			12,
		},
		{
			"POW 2 3",
			Application{Func: Application{Func: POW, Arg: Numeral(2)}, Arg: Numeral(3)},
			8,
		},
		{
			"PLUS 0 0",
			Application{Func: Application{Func: PLUS, Arg: Numeral(0)}, Arg: Numeral(0)},
			0,
		},
		{
			"MULT 0 5",
			Application{Func: Application{Func: MULT, Arg: Numeral(0)}, Arg: Numeral(5)},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Reduce(tt.term, 1000)
			got := ToInt(result)
			if got != tt.want {
				t.Errorf("%s = %d, want %d (reduced to: %s)", tt.name, got, tt.want, result)
			}
		})
	}
}

func TestNumeralMixedWithChurch(t *testing.T) {
	// Numeral + ChurchNumeral should work together
	term := Application{
		Func: Application{Func: PLUS, Arg: Numeral(3)},
		Arg:  ChurchNumeral(5),
	}
	result, _ := Reduce(term, 1000)
	got := ToInt(result)
	if got != 8 {
		t.Errorf("PLUS Numeral(3) Church(5) = %d, want 8", got)
	}
}

func TestNumeralSubstitute(t *testing.T) {
	// Numeral should pass through substitution unchanged
	n := Numeral(5)
	got := n.Substitute("x", Var{Name: "y"})
	if got != n {
		t.Errorf("Numeral.Substitute changed it: %s", got)
	}
}

func TestNumeralApplySubstitute(t *testing.T) {
	// NumeralApply should substitute in F
	na := NumeralApply{N: 3, Param: "x", F: Var{Name: "f"}}
	got := na.Substitute("f", Var{Name: "g"})
	na2, ok := got.(NumeralApply)
	if !ok {
		t.Fatalf("expected NumeralApply, got %T", got)
	}
	if na2.F.String() != "g" {
		t.Errorf("F should be g, got %s", na2.F)
	}
}

func TestNumeralDiagram(t *testing.T) {
	// Diagram of Numeral(2) should match Church 2
	got := Diagram(Numeral(2))
	want := Diagram(ChurchNumeral(2))
	if got != want {
		t.Errorf("Diagram(Numeral(2)):\n%s\nwant:\n%s", got, want)
	}
}

func TestNumeralApplyEtaConvert(t *testing.T) {
	// NumeralApply{1, f} = λx.f x → f (η-reduction)
	na := NumeralApply{N: 1, Param: "x", F: Var{Name: "f"}}
	result, converted := na.EtaConvert()
	if !converted {
		t.Fatal("should η-convert")
	}
	if v, ok := result.(Var); !ok || v.Name != "f" {
		t.Errorf("η-conversion of λx.f x = %s, want f", result)
	}
}

func BenchmarkPLUS_Church(b *testing.B) {
	term := Application{
		Func: Application{Func: PLUS, Arg: ChurchNumeral(10)},
		Arg:  ChurchNumeral(10),
	}
	for i := 0; i < b.N; i++ {
		result, _ := Reduce(term, 10000)
		ToInt(result)
	}
}

func BenchmarkPLUS_Numeral(b *testing.B) {
	term := Application{
		Func: Application{Func: PLUS, Arg: Numeral(10)},
		Arg:  Numeral(10),
	}
	for i := 0; i < b.N; i++ {
		result, _ := Reduce(term, 10000)
		ToInt(result)
	}
}

func BenchmarkMULT_Church(b *testing.B) {
	term := Application{
		Func: Application{Func: MULT, Arg: ChurchNumeral(5)},
		Arg:  ChurchNumeral(5),
	}
	for i := 0; i < b.N; i++ {
		result, _ := Reduce(term, 10000)
		ToInt(result)
	}
}

func BenchmarkMULT_Numeral(b *testing.B) {
	term := Application{
		Func: Application{Func: MULT, Arg: Numeral(5)},
		Arg:  Numeral(5),
	}
	for i := 0; i < b.N; i++ {
		result, _ := Reduce(term, 10000)
		ToInt(result)
	}
}
