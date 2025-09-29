package lambda

import (
	"fmt"
	"strings"
)

// Diagram represents a lambda diagram as a 2D grid
// Based on https://tromp.github.io/cl/diagrams.html
type Diagram struct {
	Grid   [][]rune
	Width  int
	Height int
}

// DiagramElement represents a single element in the diagram
type DiagramElement int

const (
	Empty DiagramElement = iota
	HLine              // Horizontal line (abstraction)
	VLine              // Vertical line (variable)
	Link               // Application link
	Corner             // Corner connection
	Cross              // Crossing lines
)

// NewDiagram creates a new diagram with the given dimensions
func NewDiagram(width, height int) *Diagram {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	return &Diagram{
		Grid:   grid,
		Width:  width,
		Height: height,
	}
}

// Set sets a character at the given position
func (d *Diagram) Set(row, col int, ch rune) {
	if row >= 0 && row < d.Height && col >= 0 && col < d.Width {
		d.Grid[row][col] = ch
	}
}

// Get gets a character at the given position
func (d *Diagram) Get(row, col int) rune {
	if row >= 0 && row < d.Height && col >= 0 && col < d.Width {
		return d.Grid[row][col]
	}
	return ' '
}

// ToUnicode converts the diagram to Unicode box drawing characters
func (d *Diagram) ToUnicode() string {
	var sb strings.Builder
	for i, row := range d.Grid {
		for _, ch := range row {
			sb.WriteRune(ch)
		}
		if i < len(d.Grid)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

// ToSVG converts the diagram to SVG format
func (d *Diagram) ToSVG() string {
	const cellWidth = 20
	const cellHeight = 20

	width := d.Width * cellWidth
	height := d.Height * cellHeight

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		width, height, width, height))
	sb.WriteString("\n")
	sb.WriteString(`<style>line{stroke:black;stroke-width:2;stroke-linecap:round;}text{font-family:monospace;font-size:14px;}</style>`)
	sb.WriteString("\n")

	// Draw the diagram elements
	for row := 0; row < d.Height; row++ {
		for col := 0; col < d.Width; col++ {
			ch := d.Grid[row][col]
			x := col * cellWidth + cellWidth/2
			y := row * cellHeight + cellHeight/2

			switch ch {
			case '─', '━': // Horizontal line
				x1 := col * cellWidth
				x2 := (col + 1) * cellWidth
				sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x1, y, x2, y))
				sb.WriteString("\n")
			case '│', '┃': // Vertical line
				y1 := row * cellHeight
				y2 := (row + 1) * cellHeight
				sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y1, x, y2))
				sb.WriteString("\n")
			case '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼': // Corners and intersections
				// Draw connecting lines based on the character
				drawCornerSVG(&sb, ch, x, y, cellWidth, cellHeight)
			}
		}
	}

	sb.WriteString("</svg>")
	return sb.String()
}

// drawCornerSVG draws corner and intersection characters as SVG lines
func drawCornerSVG(sb *strings.Builder, ch rune, x, y, cellWidth, cellHeight int) {
	halfW := cellWidth / 2
	halfH := cellHeight / 2

	switch ch {
	case '┌': // Top-left corner
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x+halfW, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x, y+halfH))
		sb.WriteString("\n")
	case '┐': // Top-right corner
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x, y+halfH))
		sb.WriteString("\n")
	case '└': // Bottom-left corner
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x+halfW, y))
		sb.WriteString("\n")
	case '┘': // Bottom-right corner
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x, y))
		sb.WriteString("\n")
	case '├': // Left T
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y+halfH))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x+halfW, y))
		sb.WriteString("\n")
	case '┤': // Right T
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y+halfH))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x, y))
		sb.WriteString("\n")
	case '┬': // Top T
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x+halfW, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y, x, y+halfH))
		sb.WriteString("\n")
	case '┴': // Bottom T
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x+halfW, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y))
		sb.WriteString("\n")
	case '┼': // Cross
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x-halfW, y, x+halfW, y))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x, y-halfH, x, y+halfH))
		sb.WriteString("\n")
	}
}

// DiagramContext tracks variable positions during diagram generation
type DiagramContext struct {
	VarPositions map[string][]int // Maps variable names to their column positions
	CurrentDepth int              // Current abstraction depth
	CurrentCol   int              // Current column position
}

// ToDiagram converts a lambda calculus term to a diagram
func ToDiagram(obj Object) *Diagram {
	// First pass: calculate dimensions
	width, height := calculateDimensions(obj, 0)

	// Add padding
	width += 2
	height += 2

	diagram := NewDiagram(width, height)
	ctx := &DiagramContext{
		VarPositions: make(map[string][]int),
		CurrentDepth: 1,
		CurrentCol:   1,
	}

	// Second pass: draw the diagram
	drawObject(diagram, obj, ctx, 1)

	return diagram
}

// ToDiagram method for Var
func (v Var) ToDiagram() *Diagram {
	return ToDiagram(v)
}

// ToDiagram method for Abstraction
func (a Abstraction) ToDiagram() *Diagram {
	return ToDiagram(a)
}

// ToDiagram method for Application
func (a Application) ToDiagram() *Diagram {
	return ToDiagram(a)
}

// calculateDimensions calculates the width and height needed for the diagram
func calculateDimensions(obj Object, depth int) (width, height int) {
	switch term := obj.(type) {
	case Var:
		return 2, depth + 1
	case Abstraction:
		w, h := calculateDimensions(term.Body, depth+1)
		return w + 2, max(h, depth+2)
	case Application:
		w1, h1 := calculateDimensions(term.Func, depth)
		w2, h2 := calculateDimensions(term.Arg, depth)
		return w1 + w2 + 2, max(h1, h2)
	}
	return 4, depth + 1
}

// drawObject draws a lambda calculus object into the diagram
func drawObject(d *Diagram, obj Object, ctx *DiagramContext, row int) int {
	switch term := obj.(type) {
	case Var:
		// Draw a vertical line for the variable
		col := ctx.CurrentCol
		ctx.CurrentCol += 2

		// Draw vertical line down from binding lambda
		for r := row; r < d.Height-1; r++ {
			d.Set(r, col, '│')
		}

		// Store variable position
		if ctx.VarPositions[term.Name] == nil {
			ctx.VarPositions[term.Name] = []int{}
		}
		ctx.VarPositions[term.Name] = append(ctx.VarPositions[term.Name], col)

		return col

	case Abstraction:
		// Draw horizontal line for abstraction
		startCol := ctx.CurrentCol

		// Draw the lambda line
		for c := startCol; c < startCol+4 && c < d.Width; c++ {
			d.Set(row, c, '─')
		}

		ctx.CurrentCol = startCol + 1
		ctx.CurrentDepth++

		// Draw the body
		drawObject(d, term.Body, ctx, row+1)

		ctx.CurrentDepth--

		return startCol

	case Application:
		// Draw function and argument
		funcCol := drawObject(d, term.Func, ctx, row)
		argCol := drawObject(d, term.Arg, ctx, row)

		// Draw horizontal link connecting them
		if funcCol < argCol {
			for c := funcCol; c <= argCol; c++ {
				if d.Get(row, c) == ' ' {
					d.Set(row, c, '─')
				}
			}
		}

		return funcCol
	}

	return ctx.CurrentCol
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}