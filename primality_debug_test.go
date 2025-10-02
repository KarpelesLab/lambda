package lambda

import (
	"testing"
)

// must is a test helper that panics on parse errors
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestPrimalityComponents(t *testing.T) {
	// Test: _LEQ 2 _3 should be TRUE
	t.Run("LEQ 2 3", func(t *testing.T) {
		expr := must(Parse("_LEQ _2 _3"))
		result, _ := Reduce(expr, 1000)
		got := ToBool(result)
		t.Logf("LEQ 2 3 = %v (want true)", got)
		if !got {
			t.Errorf("LEQ 2 3 = %v, want true", got)
		}
	})

	// Test: _EQ 2 _TWO should be TRUE
	t.Run("EQ 2 TWO", func(t *testing.T) {
		expr := must(Parse("_EQ _2 _TWO"))
		result, _ := Reduce(expr, 1000)
		got := ToBool(result)
		t.Logf("EQ 2 TWO = %v (want true)", got)
		if !got {
			t.Errorf("EQ 2 TWO = %v, want true", got)
		}
	})

	// Test: _OR (_EQ _2 _TWO) (_EQ _2 _3) should be TRUE
	t.Run("OR (EQ 2 TWO) (EQ 2 3)", func(t *testing.T) {
		expr := must(Parse("_OR (_EQ _2 _TWO) (_EQ _2 _3)"))
		result, _ := Reduce(expr, 1000)
		got := ToBool(result)
		t.Logf("OR (EQ 2 TWO) (EQ 2 3) = %v (want true)", got)
		if !got {
			t.Errorf("OR (EQ 2 TWO) (EQ 2 3) = %v, want true", got)
		}
	})

	// Test the full small-n branch: _IF (_LEQ _2 _3) (_OR (_EQ _2 _TWO) (_EQ _2 _3)) _FALSE
	t.Run("Full small-n branch for n=2", func(t *testing.T) {
		expr := must(Parse("_IF (_LEQ _2 _3) (_OR (_EQ _2 _TWO) (_EQ _2 _3)) _FALSE"))
		result, _ := Reduce(expr, 2000)
		got := ToBool(result)
		t.Logf("Full small-n branch for n=2 = %v (want true)", got)
		if !got {
			t.Errorf("Full small-n branch for n=2 = %v, want true", got)
		}
	})

	// Test IS_PRIME structure directly with substitution
	t.Run("IS_PRIME applied to 2", func(t *testing.T) {
		result := Application{
			Func: IS_PRIME,
			Arg:  ChurchNumeral(2),
		}

		// Reduce step by step to find where it goes wrong
		current := Term(result)
		for i := 0; i < 20; i++ {
			next, steps := Reduce(current, 1)
			if steps == 0 {
				break
			}
			t.Logf("Step %d: %s", i+1, next)
			current = next
		}

		reduced, steps := Reduce(result, 5000)
		got := ToBool(reduced)
		t.Logf("IS_PRIME(2) = %v (want true) in %d steps", got, steps)
		t.Logf("Reduced form: %s", reduced)
		if !got {
			t.Errorf("IS_PRIME(2) = %v, want true", got)
		}
	})

	// Test the exact small-n branch standalone
	t.Run("Standalone small-n branch", func(t *testing.T) {
		expr := must(Parse(`(\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _2`))
		result, _ := Reduce(expr, 2000)
		got := ToBool(result)
		t.Logf("Standalone small-n branch for n=2 = %v (want true)", got)
		if !got {
			t.Errorf("Standalone small-n branch for n=2 = %v, want true", got)
		}
	})

	// Print IS_PRIME structure
	t.Run("IS_PRIME structure", func(t *testing.T) {
		t.Logf("IS_PRIME = %s", IS_PRIME)
	})

	// Test _3 constant parsing
	t.Run("_3 constant", func(t *testing.T) {
		expr := must(Parse("_3"))
		t.Logf("_3 parsed as: %s", expr)
		t.Logf("_3 = %d", ToInt(expr))
	})

	// Test the complete IS_PRIME small-n logic
	t.Run("Complete small-n logic for n=2", func(t *testing.T) {
		// This is what IS_PRIME should do for n=2:
		// _IF (_LEQ _2 _3) (_OR (_EQ _2 _TWO) (_EQ _2 _3)) (_IF (_ISEVEN _2) _FALSE ...)
		// Which should reduce to: TRUE (since LEQ 2 3 is true, and then OR (EQ 2 2) (EQ 2 3) is true)

		// But let's build this manually to see what's happening
		expr := must(Parse(`
			(\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE) _2
		`))
		result, _ := Reduce(expr, 5000)
		got := ToBool(result)
		t.Logf("Manual IS_PRIME small-n for n=2 = %v (want true)", got)
		if !got {
			t.Errorf("Manual IS_PRIME small-n for n=2 = %v, want true", got)
		}
	})

	// Test: extract and evaluate just the condition from IS_PRIME
	t.Run("IS_PRIME condition for n=2", func(t *testing.T) {
		// IS_PRIME is: \n. IF (condition) (then) (else)
		// So IS_PRIME 2 is: IF (condition[n:=2]) (then[n:=2]) (else[n:=2])
		// Let's check what the condition evaluates to

		// Build: (\n. _LEQ n _3) _2
		condExpr := must(Parse(`(\n. _LEQ n _3) _2`))
		condResult, _ := Reduce(condExpr, 1000)
		condBool := ToBool(condResult)
		t.Logf("Condition (_LEQ 2 _3) = %v (want true)", condBool)

		// Also test the then branch
		thenExpr := must(Parse(`(\n. _OR (_EQ n _TWO) (_EQ n _3)) _2`))
		thenResult, _ := Reduce(thenExpr, 1000)
		thenBool := ToBool(thenResult)
		t.Logf("Then branch (_OR (_EQ 2 _TWO) (_EQ 2 _3)) = %v (want true)", thenBool)

		// Now test: IF TRUE (then) (else) should give (then)
		testExpr := must(Parse(`_IF _TRUE _TRUE _FALSE`))
		testResult, _ := Reduce(testExpr, 100)
		testBool := ToBool(testResult)
		t.Logf("_IF _TRUE _TRUE _FALSE = %v (want true)", testBool)
	})

	// Compare the parsed IS_PRIME first few characters
	t.Run("Compare IS_PRIME parsed vs manual", func(t *testing.T) {
		manual := must(Parse(`\n. _IF (_LEQ n _3) (_OR (_EQ n _TWO) (_EQ n _3)) _FALSE`))

		manualStr := manual.String()
		isPrimeStr := IS_PRIME.String()

		// Find where they diverge
		divergeAt := 0
		for i := 0; i < len(manualStr) && i < len(isPrimeStr); i++ {
			if manualStr[i] != isPrimeStr[i] {
				divergeAt = i
				break
			}
		}

		if divergeAt > 0 {
			t.Logf("Strings diverge at position %d", divergeAt)
			t.Logf("Manual around divergence: ...%s...", manualStr[max(0, divergeAt-20):min(len(manualStr), divergeAt+80)])
			t.Logf("IS_PRIME around divergence: ...%s...", isPrimeStr[max(0, divergeAt-20):min(len(isPrimeStr), divergeAt+80)])
		} else {
			t.Logf("Strings are identical for %d characters", min(len(manualStr), len(isPrimeStr)))
			t.Logf("Manual length: %d, IS_PRIME length: %d", len(manualStr), len(isPrimeStr))

			// What's after the common prefix in manual?
			commonLen := min(len(manualStr), len(isPrimeStr))
			if len(manualStr) > commonLen {
				t.Logf("Manual has extra: %q", manualStr[commonLen:])
			}
			if len(isPrimeStr) > commonLen {
				t.Logf("IS_PRIME has extra: %q", isPrimeStr[commonLen:])
			}

			t.Logf("Manual last 50 chars: ...%s", manualStr[max(0, len(manualStr)-50):])
			t.Logf("IS_PRIME last 50 chars: ...%s", isPrimeStr[max(0, len(isPrimeStr)-50):])
		}

		// Apply both to 2 and compare results
		manualResult, _ := Reduce(Application{Func: manual, Arg: ChurchNumeral(2)}, 5000)
		isPrimeResult, _ := Reduce(Application{Func: IS_PRIME, Arg: ChurchNumeral(2)}, 5000)

		t.Logf("Manual(2) = %v", ToBool(manualResult))
		t.Logf("IS_PRIME(2) = %v", ToBool(isPrimeResult))
	})
}