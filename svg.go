package lambda

import (
	"fmt"
	"math"
	"strings"
)

// RectKind classifies SVG rect elements in the diagram.
type RectKind int

const (
	RectLambda    RectKind = iota // Horizontal bar for abstraction
	RectVariable                  // Vertical line from binding lambda to variable use
	RectApp                       // Horizontal connector bar at bottom of application
	RectConnector                 // Vertical connector from sub-term output to app bar
)

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

func (c Color) css() string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
}

// SVGRect is a single rectangle in the Tromp diagram.
type SVGRect struct {
	ID     int
	Kind   RectKind
	Row    int
	Col    int
	Width  int
	Height int
	Color  Color
}

// SVGDiagram holds the complete diagram layout as rectangles.
type SVGDiagram struct {
	GridWidth  int
	GridHeight int
	Rects      []SVGRect
}

// SVGOptions controls rendering parameters.
type SVGOptions struct {
	CellSize   int               // Pixels per grid cell (default: 20)
	Padding    int               // Padding around diagram in pixels (default: 10)
	Background string            // Background color CSS string (default: "#000")
	LineWidth  int               // Width of lines in pixels (default: CellSize/3)
	Colors     map[int]Color     // Override colors by lambda index (pre-order)
	Saturation float64           // Default HSV saturation 0-1 (default: 0.7)
	Value      float64           // Default HSV value/brightness 0-1 (default: 1.0)
}

func (o *SVGOptions) cellSize() int {
	if o != nil && o.CellSize > 0 {
		return o.CellSize
	}
	return 20
}

func (o *SVGOptions) padding() int {
	if o != nil && o.Padding > 0 {
		return o.Padding
	}
	return 10
}

func (o *SVGOptions) background() string {
	if o != nil && o.Background != "" {
		return o.Background
	}
	return "#000"
}

func (o *SVGOptions) lineWidth() int {
	if o != nil && o.LineWidth > 0 {
		return o.LineWidth
	}
	cs := o.cellSize()
	lw := cs / 3
	if lw < 1 {
		lw = 1
	}
	return lw
}

func (o *SVGOptions) saturation() float64 {
	if o != nil && o.Saturation > 0 {
		return o.Saturation
	}
	return 0.7
}

func (o *SVGOptions) value() float64 {
	if o != nil && o.Value > 0 {
		return o.Value
	}
	return 1.0
}

func (o *SVGOptions) colorFor(lambdaIndex, totalLambdas int) Color {
	if o != nil && o.Colors != nil {
		if c, ok := o.Colors[lambdaIndex]; ok {
			return c
		}
	}
	if totalLambdas <= 0 {
		totalLambdas = 1
	}
	h := float64(lambdaIndex) / float64(totalLambdas)
	return hsvToRGB(h, o.saturation(), o.value())
}

// HSVColor creates a Color from hue (0-1), saturation (0-1), and value (0-1).
func HSVColor(h, s, v float64) Color {
	return hsvToRGB(h, s, v)
}

func hsvToRGB(h, s, v float64) Color {
	h = h - math.Floor(h) // normalize to [0,1)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := v - c
	var r, g, b float64
	switch int(h * 6) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return Color{
		R: uint8(math.Round((r + m) * 255)),
		G: uint8(math.Round((g + m) * 255)),
		B: uint8(math.Round((b + m) * 255)),
	}
}

// blendColors averages two colors for application connectors.
func blendColors(a, b Color) Color {
	return Color{
		R: uint8((int(a.R) + int(b.R)) / 2),
		G: uint8((int(a.G) + int(b.G)) / 2),
		B: uint8((int(a.B) + int(b.B)) / 2),
	}
}

// BuildSVGDiagram constructs the rect-based diagram for a term.
func BuildSVGDiagram(term Term, opts *SVGOptions) *SVGDiagram {
	db := toDeBruijn(term, nil)
	info := computeInfo(db)
	numLambdas := countLambdas(db)

	b := &svgBuilder{opts: opts, numLambdas: numLambdas}
	b.build(db, 0, 0, nil)

	return &SVGDiagram{
		GridWidth:  info.width,
		GridHeight: info.height,
		Rects:      b.rects,
	}
}

// DiagramSVG returns a static SVG string for a lambda term.
func DiagramSVG(term Term, opts *SVGOptions) string {
	d := BuildSVGDiagram(term, opts)
	return d.SVG(opts)
}

// SVG renders the diagram to an SVG string.
func (d *SVGDiagram) SVG(opts *SVGOptions) string {
	return d.renderSVG(opts, "")
}

func (d *SVGDiagram) renderSVG(opts *SVGOptions, extraCSS string) string {
	cs := opts.cellSize()
	pad := opts.padding()
	lw := opts.lineWidth()

	totalW := d.GridWidth*cs + 2*pad
	totalH := d.GridHeight*cs + 2*pad

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, totalW, totalH, totalW, totalH)
	sb.WriteByte('\n')

	if extraCSS != "" {
		fmt.Fprintf(&sb, "<style>\n%s</style>\n", extraCSS)
	}

	// Background
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="%s"/>`, totalW, totalH, opts.background())
	sb.WriteByte('\n')

	// Draw each rect
	for _, r := range d.Rects {
		x, y, w, h := rectPixels(r, cs, pad, lw)
		fmt.Fprintf(&sb, `<rect id="r%d" class="%s" x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`,
			r.ID, rectClass(r.Kind), x, y, w, h, r.Color.css())
		sb.WriteByte('\n')
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

// rectPixels computes pixel coordinates for a rect.
// Horizontal bars span full cell width, centered vertically in their cell.
// Vertical lines are centered horizontally and span from center of first cell to center of last cell,
// ensuring overlap with horizontal bars at connection points.
func rectPixels(r SVGRect, cs, pad, lw int) (x, y, w, h float64) {
	offset := float64(cs-lw) / 2
	switch r.Kind {
	case RectLambda, RectApp:
		// Horizontal bar: full cell width, narrow height, vertically centered
		x = float64(r.Col*cs + pad)
		y = float64(r.Row*cs+pad) + offset
		w = float64(r.Width * cs)
		h = float64(lw)
	case RectVariable, RectConnector:
		// Vertical line: centered in cell, from center of first row to center of last row
		x = float64(r.Col*cs+pad) + offset
		y = float64(r.Row*cs+pad) + offset
		w = float64(lw)
		h = float64((r.Height-1)*cs) + float64(lw)
	}
	return
}

func rectClass(k RectKind) string {
	switch k {
	case RectLambda:
		return "lam"
	case RectVariable:
		return "var"
	case RectApp:
		return "app"
	case RectConnector:
		return "conn"
	}
	return ""
}

// countLambdas counts the number of abstraction nodes in a de Bruijn term.
func countLambdas(t dbTerm) int {
	switch term := t.(type) {
	case dbVar:
		return 0
	case dbAbs:
		return 1 + countLambdas(term.body)
	case dbApp:
		return countLambdas(term.fun) + countLambdas(term.arg)
	}
	return 0
}

// svgBuilder accumulates rects during the recursive tree walk.
type svgBuilder struct {
	opts       *SVGOptions
	numLambdas int
	lambdaIdx  int // pre-order counter for lambda color assignment
	nextID     int
	rects      []SVGRect
}

// lambdaInfo tracks binding lambda positions and colors during recursion.
type lambdaInfo struct {
	row   int
	color Color
}

// buildResult holds the output point and color of a rendered sub-term.
type buildResult struct {
	outRow int
	outCol int
	color  Color
}

func (b *svgBuilder) build(t dbTerm, topRow, leftCol int, lambdas []lambdaInfo) buildResult {
	switch term := t.(type) {
	case dbVar:
		outRow := topRow
		outCol := leftCol
		var color Color
		if term.index < len(lambdas) {
			li := lambdas[len(lambdas)-1-term.index]
			color = li.color
			// Draw variable line from binding lambda row down to this row
			lineRow := li.row
			lineHeight := outRow - lineRow + 1
			if lineHeight > 0 {
				b.rects = append(b.rects, SVGRect{
					ID:     b.nextID,
					Kind:   RectVariable,
					Row:    lineRow,
					Col:    leftCol,
					Width:  1,
					Height: lineHeight,
					Color:  color,
				})
			}
		} else {
			// Free variable — use grey
			color = Color{128, 128, 128}
			b.rects = append(b.rects, SVGRect{
				ID:     b.nextID,
				Kind:   RectVariable,
				Row:    topRow,
				Col:    leftCol,
				Width:  1,
				Height: 1,
				Color:  color,
			})
		}
		b.nextID++
		return buildResult{outRow, outCol, color}

	case dbAbs:
		bodyInfo := computeInfo(term.body)
		lambdaColor := b.opts.colorFor(b.lambdaIdx, b.numLambdas)
		b.lambdaIdx++

		// Lambda bar rect
		b.rects = append(b.rects, SVGRect{
			ID:    b.nextID,
			Kind:  RectLambda,
			Row:   topRow,
			Col:   leftCol,
			Width: bodyInfo.width,
			Color: lambdaColor,
		})
		b.nextID++

		// Recurse into body
		newLambdas := append(lambdas, lambdaInfo{row: topRow, color: lambdaColor})
		result := b.build(term.body, topRow+1, leftCol, newLambdas)
		return buildResult{result.outRow, result.outCol, result.color}

	case dbApp:
		fInfo := computeInfo(term.fun)
		aInfo := computeInfo(term.arg)

		funResult := b.build(term.fun, topRow, leftCol, lambdas)
		argResult := b.build(term.arg, topRow, leftCol+fInfo.width+1, lambdas)

		maxH := fInfo.height
		if aInfo.height > maxH {
			maxH = aInfo.height
		}
		connRow := topRow + maxH

		appColor := blendColors(funResult.color, argResult.color)

		// Application connector bar
		barCol := funResult.outCol
		barWidth := argResult.outCol - funResult.outCol + 1
		b.rects = append(b.rects, SVGRect{
			ID:    b.nextID,
			Kind:  RectApp,
			Row:   connRow,
			Col:   barCol,
			Width: barWidth,
			Color: appColor,
		})
		b.nextID++

		// Function connector (vertical from function output down to connector bar)
		if funResult.outRow < connRow {
			b.rects = append(b.rects, SVGRect{
				ID:     b.nextID,
				Kind:   RectConnector,
				Row:    funResult.outRow,
				Col:    funResult.outCol,
				Width:  1,
				Height: connRow - funResult.outRow + 1,
				Color:  funResult.color,
			})
			b.nextID++
		}

		// Input connector (vertical from input output down to connector bar)
		if argResult.outRow < connRow {
			b.rects = append(b.rects, SVGRect{
				ID:     b.nextID,
				Kind:   RectConnector,
				Row:    argResult.outRow,
				Col:    argResult.outCol,
				Width:  1,
				Height: connRow - argResult.outRow + 1,
				Color:  argResult.color,
			})
			b.nextID++
		}

		return buildResult{connRow, funResult.outCol, appColor}
	}
	return buildResult{topRow, leftCol, Color{255, 255, 255}}
}
