package lambda

import (
	"testing"
)

func TestDiagramIdentity(t *testing.T) {
	got := Diagram(I)
	expect := "┬\n╵"
	if got != expect {
		t.Errorf("Diagram(I):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramK(t *testing.T) {
	// K = λx.λy.x — var passes through inner bar to outer
	got := Diagram(K)
	expect := "┬\n┼\n╵"
	if got != expect {
		t.Errorf("Diagram(K):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramFalse(t *testing.T) {
	// FALSE = λx.λy.y — outer bar has no var, inner bar starts var
	got := Diagram(F)
	expect := "─\n┬\n╵"
	if got != expect {
		t.Errorf("Diagram(FALSE):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramU(t *testing.T) {
	// U = λx.xx — bar, two lines, app connector at bottom
	got := Diagram(U)
	expect := "┌─┐\n│ │\n└─┘"
	if got != expect {
		t.Errorf("Diagram(U):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramOmega(t *testing.T) {
	// OMEGA = (λx.xx)(λx.xx) — two sub-diagrams, outer connector at bottom
	got := Diagram(OMEGA)
	expect := "┌─┐ ┌─┐\n│ │ │ │\n├─┘ ├─┘\n└───┘"
	if got != expect {
		t.Errorf("Diagram(OMEGA):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramChurch2(t *testing.T) {
	// 2 = λf.λx.f(fx) — converges to single output at bottom
	got := Diagram(ChurchNumeral(2))
	expect := "┌─┬─╴\n├─┼─┐\n│ │ │\n│ ├─┘\n└─┘"
	if got != expect {
		t.Errorf("Diagram(2):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramChurch3(t *testing.T) {
	got := Diagram(ChurchNumeral(3))
	expect := "┌─┬─┬─╴\n├─┼─┼─┐\n│ │ │ │\n│ │ ├─┘\n│ ├─┘\n└─┘"
	if got != expect {
		t.Errorf("Diagram(3):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramS(t *testing.T) {
	got := Diagram(S)
	expect := "┌─────╴\n├───┬─╴\n├─┬─┼─┐\n│ │ │ │\n├─┘ ├─┘\n└───┘"
	if got != expect {
		t.Errorf("Diagram(S):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramY(t *testing.T) {
	got := Diagram(Y)
	expect := "┌─────┬───╴\n├─┬─┐ ├─┬─┐\n│ │ │ │ │ │\n│ ├─┘ │ ├─┘\n├─┘   ├─┘\n└─────┘"
	if got != expect {
		t.Errorf("Diagram(Y):\ngot:\n%s\nexpect:\n%s", got, expect)
	}
}

func TestDiagramSingleOutput(t *testing.T) {
	// All closed terms should have a single connected structure at the bottom row.
	// Verify by checking that the bottom row has no gaps between drawn characters.
	terms := []struct {
		name string
		term Term
	}{
		{"I", I},
		{"K", K},
		{"S", S},
		{"U", U},
		{"OMEGA", OMEGA},
		{"Y", Y},
		{"0", ChurchNumeral(0)},
		{"1", ChurchNumeral(1)},
		{"2", ChurchNumeral(2)},
		{"3", ChurchNumeral(3)},
		{"4", ChurchNumeral(4)},
	}
	for _, tt := range terms {
		t.Run(tt.name, func(t *testing.T) {
			diagram := Diagram(tt.term)
			lines := splitLines(diagram)
			if len(lines) == 0 {
				t.Fatal("empty diagram")
			}
			bottom := lines[len(lines)-1]
			// Bottom line should have no internal spaces (single connected piece)
			inContent := false
			for _, ch := range bottom {
				if ch != ' ' {
					inContent = true
				} else if inContent {
					t.Errorf("bottom line has gap: %q", bottom)
					break
				}
			}
		})
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
