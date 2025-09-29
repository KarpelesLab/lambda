package lambda

import (
	"strings"
	"testing"
)

func TestDiagramIdentity(t *testing.T) {
	// λx.x (identity function)
	identity := Abstraction{
		Param: "x",
		Body:  Var{Name: "x"},
	}

	diagram := identity.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("Identity function diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramK(t *testing.T) {
	// λx.λy.x (K combinator / TRUE)
	k := Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body:  Var{Name: "x"},
		},
	}

	diagram := k.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("K combinator diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramChurchNumeral(t *testing.T) {
	// Church numeral 2: λf.λx.f (f x)
	two := ChurchNumeral(2)

	diagram := ToDiagram(two)
	unicode := diagram.ToUnicode()

	t.Logf("Church numeral 2 diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramSVG(t *testing.T) {
	// Test SVG generation
	identity := Abstraction{
		Param: "x",
		Body:  Var{Name: "x"},
	}

	diagram := identity.ToDiagram()
	svg := diagram.ToSVG()

	// Check that SVG contains expected elements
	if !strings.Contains(svg, "<svg") {
		t.Error("SVG output doesn't contain <svg> tag")
	}

	if !strings.Contains(svg, "</svg>") {
		t.Error("SVG output doesn't contain closing </svg> tag")
	}

	if !strings.Contains(svg, "<line") {
		t.Error("SVG output doesn't contain any <line> elements")
	}

	t.Logf("SVG output length: %d bytes", len(svg))
}

func TestDiagramApplication(t *testing.T) {
	// (λx.x) y
	app := Application{
		Func: Abstraction{
			Param: "x",
			Body:  Var{Name: "x"},
		},
		Arg: Var{Name: "y"},
	}

	diagram := app.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("Application diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramComplex(t *testing.T) {
	// More complex: λf.λx.f (f x) - Church numeral 2
	term := Abstraction{
		Param: "f",
		Body: Abstraction{
			Param: "x",
			Body: Application{
				Func: Var{Name: "f"},
				Arg: Application{
					Func: Var{Name: "f"},
					Arg:  Var{Name: "x"},
				},
			},
		},
	}

	diagram := term.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("Church numeral 2 (explicit) diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramTRUE(t *testing.T) {
	diagram := TRUE.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("TRUE diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}

func TestDiagramFALSE(t *testing.T) {
	diagram := FALSE.ToDiagram()
	unicode := diagram.ToUnicode()

	t.Logf("FALSE diagram:\n%s", unicode)

	if diagram.Width == 0 || diagram.Height == 0 {
		t.Error("Diagram has zero dimensions")
	}
}