package lambda

import (
	"fmt"
	"testing"
)

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		want  interface{} // can be int or bool
		steps int
	}{
		// Basic constants
		{"Zero", "_0", 0, 100},
		{"One", "_1", 1, 100},
		{"Two", "_2", 2, 100},
		{"Three", "_3", 3, 100},
		{"TWO constant", "_TWO", 2, 100},

		// Boolean operations
		{"TRUE", "_TRUE", true, 100},
		{"FALSE", "_FALSE", false, 100},
		{"NOT TRUE", "_NOT _TRUE", false, 100},
		{"NOT FALSE", "_NOT _FALSE", true, 100},
		{"AND TRUE TRUE", "_AND _TRUE _TRUE", true, 100},
		{"AND TRUE FALSE", "_AND _TRUE _FALSE", false, 100},
		{"OR TRUE FALSE", "_OR _TRUE _FALSE", true, 100},
		{"OR FALSE FALSE", "_OR _FALSE _FALSE", false, 100},

		// IF expressions
		{"IF TRUE 1 0", "_IF _TRUE _1 _0", 1, 100},
		{"IF FALSE 1 0", "_IF _FALSE _1 _0", 0, 100},
		{"IF (EQ 2 2) TRUE FALSE", "_IF (_EQ _2 _2) _TRUE _FALSE", true, 1000},

		// Arithmetic
		{"SUCC 0", "_SUCC _0", 1, 100},
		{"SUCC 2", "_SUCC _2", 3, 100},
		{"PLUS 2 3", "_PLUS _2 _3", 5, 100},
		{"ADD 2 3", "_ADD _2 _3", 5, 100},
		{"MULT 2 3", "_MULT _2 _3", 6, 500},
		{"MULT 3 4", "_MULT _3 _4", 12, 1000},
		{"POW 2 3", "_POW _2 _3", 8, 1000},

		// Comparisons
		{"LEQ 2 3", "_LEQ _2 _3", true, 1000},
		{"LEQ 3 2", "_LEQ _3 _2", false, 1000},
		{"LEQ 2 2", "_LEQ _2 _2", true, 1000},
		{"LT 2 3", "_LT _2 _3", true, 1000},
		{"LT 3 3", "_LT _3 _3", false, 1000},
		{"LT 3 2", "_LT _3 _2", false, 1000},
		{"EQ 2 2", "_EQ _2 _2", true, 1000},
		{"EQ 2 3", "_EQ _2 _3", false, 1000},
		{"EQ 5 5", "_EQ _5 _5", true, 2000},

		// With variable substitution
		{"(λn. SUCC n) 2", "(\\n. _SUCC n) _2", 3, 100},
		{"(λn. PLUS n 1) 2", "(\\n. _PLUS n _1) _2", 3, 100},
		{"(λn. EQ n 2) 2", "(\\n. _EQ n _2) _2", true, 1000},
		{"(λn. LEQ n 3) 2", "(\\n. _LEQ n _3) _2", true, 1000},

		// Nested IF with variables
		{"(λn. IF (LEQ n 3) TRUE FALSE) 2", "(\\n. _IF (_LEQ n _3) _TRUE _FALSE) _2", true, 2000},
		{"(λn. IF (LEQ n 3) TRUE FALSE) 5", "(\\n. _IF (_LEQ n _3) _TRUE _FALSE) _5", false, 2000},

		// OR with variables
		{"(λn. OR (EQ n 2) (EQ n 3)) 2", "(\\n. _OR (_EQ n _2) (_EQ n _3)) _2", true, 2000},
		{"(λn. OR (EQ n 2) (EQ n 3)) 3", "(\\n. _OR (_EQ n _2) (_EQ n _3)) _3", true, 2000},
		{"(λn. OR (EQ n 2) (EQ n 3)) 4", "(\\n. _OR (_EQ n _2) (_EQ n _3)) _4", false, 2000},

		// Complex nested expression matching IS_PRIME small-n branch
		{"Small-n branch for 2", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _2", true, 3000},
		{"Small-n branch for 3", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _3", true, 3000},
		{"Small-n branch for 4", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _4", false, 3000},
		{"Small-n branch for 1", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _1", false, 3000},

		// Test ISEVEN
		{"ISEVEN 2", "_ISEVEN _2", true, 1000},
		{"ISEVEN 3", "_ISEVEN _3", false, 1000},
		{"ISEVEN 4", "_ISEVEN _4", true, 2000},

		// Nested IF-ISEVEN pattern
		{"IF ISEVEN-2 FALSE TRUE", "_IF (_ISEVEN _2) _FALSE _TRUE", false, 2000},
		{"IF ISEVEN-3 FALSE TRUE", "_IF (_ISEVEN _3) _FALSE _TRUE", true, 2000},

		// Two-level nested IF like IS_PRIME structure
		{"Two-level IF for n=2", "(\\n. _IF (_LEQ n _3) _TRUE (_IF (_ISEVEN n) _FALSE _TRUE)) _2", true, 3000},
		{"Two-level IF for n=4", "(\\n. _IF (_LEQ n _3) _TRUE (_IF (_ISEVEN n) _FALSE _TRUE)) _4", false, 3000},
		{"Two-level IF for n=5", "(\\n. _IF (_LEQ n _3) _TRUE (_IF (_ISEVEN n) _FALSE _TRUE)) _5", true, 3000},

		// Building up to IS_PRIME structure progressively
		{"IS_PRIME-like step 1: outer IF", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _2", true, 3000},
		{"IS_PRIME-like step 2: with inner IF", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) (_IF (_ISEVEN n) _FALSE _TRUE)) _2", true, 3000},
		{"IS_PRIME-like step 3: with LET", "(\\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) (_IF (_ISEVEN n) _FALSE ((\\x.\\f. f x) _ZERO (\\dummy. _TRUE)))) _2", true, 5000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result, steps := Reduce(expr, tt.steps)

			switch want := tt.want.(type) {
			case int:
				got := ToInt(result)
				if got != want {
					t.Errorf("%s = %d, want %d (in %d steps)", tt.expr, got, want, steps)
				} else {
					t.Logf("%s = %d ✓ (in %d steps)", tt.name, got, steps)
				}
			case bool:
				got := ToBool(result)
				if got != want {
					t.Errorf("%s = %v, want %v (in %d steps)", tt.expr, got, want, steps)
					t.Logf("Reduced to: %s", result)
				} else {
					t.Logf("%s = %v ✓ (in %d steps)", tt.name, got, steps)
				}
			}
		})
	}
}

func TestISPRIMEDirect(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{2, true},
		{3, true},
		{4, false},
		// IS_PRIME(5) and higher are too computationally expensive for automated tests
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: IS_PRIME,
				Arg:  ChurchNumeral(tt.n),
			}
			reduced, steps := Reduce(result, 50000)
			got := ToBool(reduced)

			if got != tt.want {
				t.Errorf("IS_PRIME(%d) = %v, want %v (in %d steps)", tt.n, got, tt.want, steps)
				t.Logf("Reduced to: %s", reduced)

				// Debug: check what the condition evaluates to
				condExpr, _ := Parse("(\\n. _LEQ n _3) " + fmt.Sprintf("_%d", tt.n))
				condResult, _ := Reduce(condExpr, 1000)
				t.Logf("Condition (LEQ %d 3) = %v", tt.n, ToBool(condResult))
			} else {
				t.Logf("IS_PRIME(%d) = %v ✓ (in %d steps)", tt.n, got, steps)
			}
		})
	}
}