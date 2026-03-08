package lambda

import (
	"strings"
	"testing"
)

func TestSVGDiagramIdentity(t *testing.T) {
	svg := DiagramSVG(I, nil)
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing </svg> closing tag")
	}
	// Identity has 1 lambda bar + 1 variable line = 2 rects (besides background)
	d := BuildSVGDiagram(I, nil)
	if len(d.Rects) != 2 {
		t.Errorf("expected 2 rects for I, got %d", len(d.Rects))
	}
	// Check kinds
	hasLambda, hasVar := false, false
	for _, r := range d.Rects {
		switch r.Kind {
		case RectLambda:
			hasLambda = true
		case RectVariable:
			hasVar = true
		}
	}
	if !hasLambda {
		t.Error("missing lambda rect")
	}
	if !hasVar {
		t.Error("missing variable rect")
	}
}

func TestSVGDiagramOmega(t *testing.T) {
	d := BuildSVGDiagram(OMEGA, nil)
	// OMEGA = (λx.xx)(λx.xx)
	// Each λx.xx: 1 lambda bar + 2 variable lines + 1 app bar + connectors
	// Plus outer app bar + connectors
	if len(d.Rects) == 0 {
		t.Error("no rects generated")
	}
	// Should have exactly 2 lambda bars
	lambdaCount := 0
	for _, r := range d.Rects {
		if r.Kind == RectLambda {
			lambdaCount++
		}
	}
	if lambdaCount != 2 {
		t.Errorf("expected 2 lambda rects, got %d", lambdaCount)
	}
}

func TestSVGDiagramGridSize(t *testing.T) {
	// Grid sizes should match the text diagram computeInfo
	tests := []struct {
		name string
		term Term
		w, h int
	}{
		{"I", I, 1, 2},
		{"K", K, 1, 3},
		{"U", U, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := BuildSVGDiagram(tt.term, nil)
			if d.GridWidth != tt.w || d.GridHeight != tt.h {
				t.Errorf("grid size %dx%d, want %dx%d", d.GridWidth, d.GridHeight, tt.w, tt.h)
			}
		})
	}
}

func TestSVGCustomOptions(t *testing.T) {
	opts := &SVGOptions{
		CellSize:   30,
		Padding:    20,
		Background: "#fff",
		LineWidth:  5,
	}
	svg := DiagramSVG(I, opts)
	if !strings.Contains(svg, `fill="#fff"`) {
		t.Error("custom background not applied")
	}
	// viewBox should reflect custom sizes: (1*30 + 2*20) x (2*30 + 2*20) = 70 x 100
	if !strings.Contains(svg, `viewBox="0 0 70 100"`) {
		t.Errorf("unexpected viewBox in:\n%s", svg)
	}
}

func TestSVGCustomColors(t *testing.T) {
	red := Color{255, 0, 0}
	opts := &SVGOptions{
		Colors: map[int]Color{0: red},
	}
	d := BuildSVGDiagram(I, opts)
	// The single lambda should have the custom color
	for _, r := range d.Rects {
		if r.Kind == RectLambda {
			if r.Color != red {
				t.Errorf("lambda color = %v, want red", r.Color)
			}
		}
	}
}

func TestSVGRectClasses(t *testing.T) {
	svg := DiagramSVG(U, nil) // U has all rect types
	for _, class := range []string{"lam", "var", "app"} {
		if !strings.Contains(svg, `class="`+class+`"`) {
			t.Errorf("missing class %q in SVG", class)
		}
	}
}

func TestSVGDistinctColors(t *testing.T) {
	// Church numeral 3 has 2 lambdas — they should get different colors
	d := BuildSVGDiagram(ChurchNumeral(3), nil)
	var lambdaColors []Color
	for _, r := range d.Rects {
		if r.Kind == RectLambda {
			lambdaColors = append(lambdaColors, r.Color)
		}
	}
	if len(lambdaColors) != 2 {
		t.Fatalf("expected 2 lambda colors, got %d", len(lambdaColors))
	}
	if lambdaColors[0] == lambdaColors[1] {
		t.Error("both lambdas got the same color")
	}
}

func TestSVGBlendColors(t *testing.T) {
	a := Color{100, 200, 50}
	b := Color{200, 100, 150}
	got := blendColors(a, b)
	want := Color{150, 150, 100}
	if got != want {
		t.Errorf("blendColors(%v, %v) = %v, want %v", a, b, got, want)
	}
}

func TestHSVColor(t *testing.T) {
	// Red at full saturation and value
	c := HSVColor(0, 1.0, 1.0)
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("HSVColor(0,1,1) = %v, want {255,0,0}", c)
	}
	// Green
	c = HSVColor(1.0/3, 1.0, 1.0)
	if c.G != 255 {
		t.Errorf("HSVColor(1/3,1,1) green = %d, want 255", c.G)
	}
}

func TestSVGWellFormed(t *testing.T) {
	// Test several terms produce well-formed SVG
	terms := []struct {
		name string
		term Term
	}{
		{"I", I},
		{"K", K},
		{"S", S},
		{"OMEGA", OMEGA},
		{"Y", Y},
		{"Church0", ChurchNumeral(0)},
		{"Church3", ChurchNumeral(3)},
	}
	for _, tt := range terms {
		t.Run(tt.name, func(t *testing.T) {
			svg := DiagramSVG(tt.term, nil)
			if !strings.HasPrefix(svg, "<svg") {
				t.Error("doesn't start with <svg")
			}
			if !strings.HasSuffix(svg, "</svg>\n") {
				t.Error("doesn't end with </svg>")
			}
			// Every rect should have positive dimensions
			d := BuildSVGDiagram(tt.term, nil)
			for _, r := range d.Rects {
				if r.Width <= 0 && r.Height <= 0 {
					t.Errorf("rect %d has zero/negative dimensions: %dx%d", r.ID, r.Width, r.Height)
				}
			}
		})
	}
}
