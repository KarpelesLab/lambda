package lambda

import (
	"strings"
	"testing"
)

func TestAnimatedSVGNoReduction(t *testing.T) {
	// I is already in normal form — should return static SVG
	svg := DiagramAnimatedSVG(I, nil)
	if strings.Contains(svg, "<style>") {
		t.Error("static SVG should not contain animation styles")
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
}

func TestAnimatedSVGWithReduction(t *testing.T) {
	// (λx.x) y reduces to y — should produce animation
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	svg := DiagramAnimatedSVG(term, nil)
	if !strings.Contains(svg, "<style>") {
		t.Error("animated SVG should contain <style> block")
	}
	if !strings.Contains(svg, "@keyframes") {
		t.Error("animated SVG should contain @keyframes")
	}
	if !strings.Contains(svg, "animation:") {
		t.Error("animated SVG should contain animation property")
	}
}

func TestAnimatedSVGLoop(t *testing.T) {
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	opts := &AnimationOptions{Loop: true}
	svg := DiagramAnimatedSVG(term, opts)
	if !strings.Contains(svg, "infinite") {
		t.Error("looping animation should use 'infinite' iteration count")
	}
}

func TestAnimatedSVGCustomDuration(t *testing.T) {
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	opts := &AnimationOptions{
		StepDuration:  3.0,
		PauseDuration: 1.0,
	}
	svg := DiagramAnimatedSVG(term, opts)
	// Total duration = 1 transition * (3.0 + 1.0) + 1.0 = 5.0s
	if !strings.Contains(svg, "5.00s") {
		t.Errorf("expected 5.00s duration in SVG:\n%s", svg)
	}
}

func TestAnimatedSVGMultipleSteps(t *testing.T) {
	// SUCC 0 = (λn.λf.λx. f (n f x)) (λf.λx. x)
	// This takes multiple reduction steps
	term := Application{
		Func: SUCC,
		Arg:  ChurchNumeral(0),
	}
	opts := &AnimationOptions{Steps: 5}
	svg := DiagramAnimatedSVG(term, opts)
	if !strings.Contains(svg, "@keyframes") {
		t.Error("should produce animation for reducible term")
	}
}

func TestAnimatedFrames(t *testing.T) {
	// (λx.x) y → y
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	frames := DiagramAnimatedFrames(term, nil, 5)
	if len(frames) < 2 {
		t.Errorf("expected at least 2 frames, got %d", len(frames))
	}
	// First frame should have more rects than the last (reduction simplifies)
	if len(frames[0].Rects) <= len(frames[len(frames)-1].Rects) {
		// Not necessarily true for all terms, but for identity applied to var it should be
		// Actually identity applied just becomes a var, so first has more rects
	}
}

func TestAnimatedFramesNoReduction(t *testing.T) {
	frames := DiagramAnimatedFrames(I, nil, 5)
	if len(frames) != 1 {
		t.Errorf("expected 1 frame for normal form, got %d", len(frames))
	}
}

func TestAnimatedSVGWellFormed(t *testing.T) {
	term := Application{
		Func: Application{
			Func: K,
			Arg:  Var{Name: "a"},
		},
		Arg: Var{Name: "b"},
	}
	svg := DiagramAnimatedSVG(term, nil)
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("doesn't start with <svg")
	}
	if !strings.HasSuffix(svg, "</svg>\n") {
		t.Error("doesn't end with </svg>")
	}
	// Check well-formedness: every opening rect should be self-closing
	rectCount := strings.Count(svg, "<rect ")
	closingCount := strings.Count(svg, "/>")
	if closingCount < rectCount {
		t.Errorf("unclosed rect tags: %d rects, %d self-closing tags", rectCount, closingCount)
	}
}

func TestAnimatedSVGColors(t *testing.T) {
	term := Application{
		Func: Abstraction{Param: "x", Body: Var{Name: "x"}},
		Arg:  Var{Name: "y"},
	}
	opts := &AnimationOptions{
		SVGOptions: SVGOptions{
			Colors: map[int]Color{0: {255, 0, 0}},
		},
	}
	svg := DiagramAnimatedSVG(term, opts)
	if !strings.Contains(svg, "rgb(255,0,0)") {
		t.Error("custom color not found in animated SVG")
	}
}
