package lambda

import (
	"fmt"
	"strings"
)

// AnimationOptions controls the animated SVG output.
type AnimationOptions struct {
	SVGOptions                // Embedded rendering options
	Steps      int           // Number of beta-reduction steps to animate (default: 10)
	StepDuration float64     // Duration of each step in seconds (default: 2.0)
	PauseDuration float64    // Pause between steps in seconds (default: 0.5)
	Loop       bool          // Whether the animation loops (default: false)
}

func (o *AnimationOptions) steps() int {
	if o != nil && o.Steps > 0 {
		return o.Steps
	}
	return 10
}

func (o *AnimationOptions) stepDuration() float64 {
	if o != nil && o.StepDuration > 0 {
		return o.StepDuration
	}
	return 2.0
}

func (o *AnimationOptions) pauseDuration() float64 {
	if o != nil && o.PauseDuration > 0 {
		return o.PauseDuration
	}
	return 0.5
}

// animFrame holds one frame of the animation: the diagram at a particular reduction step.
type animFrame struct {
	diagram *SVGDiagram
	term    Term
}

// DiagramAnimatedSVG returns an SVG string with CSS animations showing beta-reduction steps.
// Each reduction step smoothly morphs rectangles from one position/size to the next.
func DiagramAnimatedSVG(term Term, opts *AnimationOptions) string {
	var svgOpts *SVGOptions
	if opts != nil {
		svgOpts = &opts.SVGOptions
	}

	maxSteps := 10
	if opts != nil {
		maxSteps = opts.steps()
	}

	// Generate frames by reducing the term
	frames := []animFrame{{
		diagram: BuildSVGDiagram(term, svgOpts),
		term:    term,
	}}

	current := term
	for i := 0; i < maxSteps; i++ {
		reduced, didReduce := current.BetaReduce()
		if !didReduce {
			break
		}
		current = reduced
		frames = append(frames, animFrame{
			diagram: BuildSVGDiagram(current, svgOpts),
			term:    current,
		})
	}

	if len(frames) <= 1 {
		// No reductions possible — return static SVG
		return frames[0].diagram.SVG(svgOpts)
	}

	return renderAnimatedSVG(frames, opts)
}

func renderAnimatedSVG(frames []animFrame, opts *AnimationOptions) string {
	var svgOpts *SVGOptions
	if opts != nil {
		svgOpts = &opts.SVGOptions
	}

	cs := svgOpts.cellSize()
	pad := svgOpts.padding()
	lw := svgOpts.lineWidth()

	stepDur := opts.stepDuration()
	pauseDur := opts.pauseDuration()
	loop := opts != nil && opts.Loop

	// Find max dimensions across all frames
	maxW, maxH := 0, 0
	for _, f := range frames {
		if f.diagram.GridWidth > maxW {
			maxW = f.diagram.GridWidth
		}
		if f.diagram.GridHeight > maxH {
			maxH = f.diagram.GridHeight
		}
	}

	totalW := maxW*cs + 2*pad
	totalH := maxH*cs + 2*pad

	// Timeline layout per frame:
	//   Frame 0: [hold pauseDur] → [transition stepDur to frame 1] →
	//   Frame 1: [hold pauseDur] → [transition stepDur to frame 2] →
	//   ...
	//   Frame N: [hold pauseDur]
	// Total = N * (pauseDur + stepDur) + pauseDur
	numTransitions := len(frames) - 1
	totalDur := float64(numTransitions)*(stepDur+pauseDur) + pauseDur

	// Maximum number of rects in any frame determines our rect pool size
	maxRects := 0
	for _, f := range frames {
		if len(f.diagram.Rects) > maxRects {
			maxRects = len(f.diagram.Rects)
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, totalW, totalH, totalW, totalH)
	sb.WriteByte('\n')

	// Generate CSS animations
	sb.WriteString("<style>\n")

	iterCount := "1"
	fillMode := "forwards"
	if loop {
		iterCount = "infinite"
		fillMode = "none"
	}

	for i := 0; i < maxRects; i++ {
		animName := fmt.Sprintf("a%d", i)
		fmt.Fprintf(&sb, "#r%d { animation: %s %.2fs ease-in-out %s; animation-fill-mode: %s; }\n",
			i, animName, totalDur, iterCount, fillMode)

		// Build keyframes: each frame holds, then transitions to the next.
		// Frame fi hold: [fi*(pause+step)] to [fi*(pause+step)+pause]
		// Transition to fi+1: [fi*(pause+step)+pause] to [(fi+1)*(pause+step)]
		fmt.Fprintf(&sb, "@keyframes %s {\n", animName)

		for fi, f := range frames {
			holdStart := float64(fi) * (stepDur + pauseDur)
			holdEnd := holdStart + pauseDur

			holdStartPct := holdStart / totalDur * 100
			holdEndPct := holdEnd / totalDur * 100

			if holdStartPct > 100 {
				holdStartPct = 100
			}
			if holdEndPct > 100 {
				holdEndPct = 100
			}

			props := rectAnimProps(i, f.diagram, cs, pad, lw)

			// Emit keyframe at hold start (arrival at this frame's state)
			// and at hold end (right before transition begins)
			fmt.Fprintf(&sb, "  %.2f%% { %s }\n", holdStartPct, props)
			if fi < len(frames)-1 {
				// Hold end — same state, transition starts after this point
				fmt.Fprintf(&sb, "  %.2f%% { %s }\n", holdEndPct, props)
			}
			// The browser interpolates from holdEnd% to next frame's holdStart%
		}
		sb.WriteString("}\n")
	}

	sb.WriteString("</style>\n")

	// Background
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="%s"/>`, totalW, totalH, svgOpts.background())
	sb.WriteByte('\n')

	// Draw rects from the first frame as initial state
	for i := 0; i < maxRects; i++ {
		if i < len(frames[0].diagram.Rects) {
			r := frames[0].diagram.Rects[i]
			x, y, w, h := rectPixels(r, cs, pad, lw)
			fmt.Fprintf(&sb, `<rect id="r%d" class="%s" x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`,
				i, rectClass(r.Kind), x, y, w, h, r.Color.css())
		} else {
			fmt.Fprintf(&sb, `<rect id="r%d" x="0" y="0" width="0" height="0" fill="transparent" opacity="0"/>`, i)
		}
		sb.WriteByte('\n')
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

// rectAnimProps returns a CSS property string for rect i in the given diagram.
func rectAnimProps(i int, d *SVGDiagram, cs, pad, lw int) string {
	if i < len(d.Rects) {
		r := d.Rects[i]
		x, y, w, h := rectPixels(r, cs, pad, lw)
		return fmt.Sprintf("x: %.1f; y: %.1f; width: %.1f; height: %.1f; fill: %s; opacity: 1",
			x, y, w, h, r.Color.css())
	}
	return "opacity: 0; width: 0; height: 0"
}

// DiagramAnimatedFrames returns individual SVG diagrams for each beta-reduction step.
// This is useful for generating frame-by-frame animations outside of CSS.
func DiagramAnimatedFrames(term Term, opts *SVGOptions, maxSteps int) []*SVGDiagram {
	if maxSteps <= 0 {
		maxSteps = 10
	}

	frames := []*SVGDiagram{BuildSVGDiagram(term, opts)}

	current := term
	for i := 0; i < maxSteps; i++ {
		reduced, didReduce := current.BetaReduce()
		if !didReduce {
			break
		}
		current = reduced
		frames = append(frames, BuildSVGDiagram(current, opts))
	}

	return frames
}
