package lambda

import "strings"

// Diagram returns a Tromp-style lambda diagram rendered with Unicode box-drawing characters.
// See https://tromp.github.io/cl/diagrams.html for the visual notation:
//   - Abstractions are horizontal bars at the top spanning all variables in scope
//   - Variables are vertical lines dropping from their binding lambda bar
//   - Applications are horizontal connectors at the bottom linking function and argument outputs
func Diagram(term Term) string {
	db := toDeBruijn(term, nil)
	info := computeInfo(db)
	grid := newGrid(info.width, info.height)

	drawTerm(db, 0, 0, nil, grid)
	return grid.String()
}

// De Bruijn representation
type dbTerm interface{ dbTag() }
type dbVar struct{ index int }
type dbAbs struct{ body dbTerm }
type dbApp struct{ fun, arg dbTerm }

func (dbVar) dbTag() {}
func (dbAbs) dbTag() {}
func (dbApp) dbTag() {}

func toDeBruijn(t Term, env []string) dbTerm {
	if ls, ok := t.(*LazyScript); ok {
		t = ls.parse()
	}
	switch term := t.(type) {
	case Var:
		for i := len(env) - 1; i >= 0; i-- {
			if env[i] == term.Name {
				return dbVar{index: len(env) - 1 - i}
			}
		}
		return dbVar{index: len(env)}
	case Abstraction:
		return dbAbs{body: toDeBruijn(term.Body, append(env, term.Param))}
	case Application:
		return dbApp{
			fun: toDeBruijn(term.Func, env),
			arg: toDeBruijn(term.Arg, env),
		}
	}
	return dbVar{index: 0}
}

// termInfo holds precomputed size and root column for a term.
type termInfo struct {
	width   int
	height  int
	rootCol int // column of the output wire (relative to term's left edge)
}

func computeInfo(t dbTerm) termInfo {
	switch term := t.(type) {
	case dbVar:
		return termInfo{width: 1, height: 1, rootCol: 0}
	case dbAbs:
		body := computeInfo(term.body)
		return termInfo{width: body.width, height: 1 + body.height, rootCol: body.rootCol}
	case dbApp:
		f := computeInfo(term.fun)
		a := computeInfo(term.arg)
		h := f.height
		if a.height > h {
			h = a.height
		}
		return termInfo{
			width:   f.width + 1 + a.width,
			height:  h + 1,
			rootCol: f.rootCol,
		}
	}
	return termInfo{}
}

// grid tracks directional connections at each cell for proper box-drawing.
type grid struct {
	w, h  int
	up    [][]bool
	down  [][]bool
	left  [][]bool
	right [][]bool
}

func newGrid(w, h int) *grid {
	g := &grid{w: w, h: h}
	g.up = makeBoolGrid(w, h)
	g.down = makeBoolGrid(w, h)
	g.left = makeBoolGrid(w, h)
	g.right = makeBoolGrid(w, h)
	return g
}

func makeBoolGrid(w, h int) [][]bool {
	g := make([][]bool, h)
	for i := range g {
		g[i] = make([]bool, w)
	}
	return g
}

func (g *grid) hLine(row, c1, c2 int) {
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	if c1 == c2 {
		// Single-cell bar: mark as horizontal so it renders as ─
		g.left[row][c1] = true
		g.right[row][c1] = true
		return
	}
	for c := c1; c <= c2; c++ {
		if c > c1 {
			g.left[row][c] = true
		}
		if c < c2 {
			g.right[row][c] = true
		}
	}
}

func (g *grid) vLine(col, r1, r2 int) {
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	for r := r1; r <= r2; r++ {
		if r > r1 {
			g.up[r][col] = true
		}
		if r < r2 {
			g.down[r][col] = true
		}
	}
}

func (g *grid) String() string {
	var sb strings.Builder
	for r := 0; r < g.h; r++ {
		if r > 0 {
			sb.WriteByte('\n')
		}
		line := make([]rune, g.w)
		for c := 0; c < g.w; c++ {
			line[c] = boxChar(g.up[r][c], g.down[r][c], g.left[r][c], g.right[r][c])
		}
		sb.WriteString(strings.TrimRight(string(line), " "))
	}
	return sb.String()
}

func boxChar(up, down, left, right bool) rune {
	switch {
	case up && down && left && right:
		return '┼'
	case up && down && right:
		return '├'
	case up && down && left:
		return '┤'
	case up && left && right:
		return '┴'
	case down && left && right:
		return '┬'
	case up && down:
		return '│'
	case left && right:
		return '─'
	case up && right:
		return '└'
	case up && left:
		return '┘'
	case down && right:
		return '┌'
	case down && left:
		return '┐'
	case up:
		return '╵'
	case down:
		return '╷'
	case left:
		return '╴'
	case right:
		return '╶'
	default:
		return ' '
	}
}

// drawTerm recursively draws the term onto the grid.
// topRow/leftCol: absolute position of this term's top-left corner.
// lambdaRows: maps de Bruijn indices to absolute rows of binding lambda bars.
// Returns the absolute row and column of this term's output point.
func drawTerm(t dbTerm, topRow, leftCol int, lambdaRows []int, g *grid) (outRow, outCol int) {
	switch term := t.(type) {
	case dbVar:
		// Variable: vertical line from binding lambda down to this position.
		outRow = topRow
		outCol = leftCol
		if term.index < len(lambdaRows) {
			bindRow := lambdaRows[len(lambdaRows)-1-term.index]
			g.vLine(leftCol, bindRow, topRow)
		}
		return

	case dbAbs:
		// Lambda bar at topRow, body below.
		bodyInfo := computeInfo(term.body)
		g.hLine(topRow, leftCol, leftCol+bodyInfo.width-1)
		newStack := append(lambdaRows, topRow)
		return drawTerm(term.body, topRow+1, leftCol, newStack, g)

	case dbApp:
		fInfo := computeInfo(term.fun)
		aInfo := computeInfo(term.arg)

		// Sub-terms are top-aligned, starting at topRow.
		funOutRow, funOutCol := drawTerm(term.fun, topRow, leftCol, lambdaRows, g)
		argOutRow, argOutCol := drawTerm(term.arg, topRow, leftCol+fInfo.width+1, lambdaRows, g)

		// App connector bar at the bottom row.
		maxH := fInfo.height
		if aInfo.height > maxH {
			maxH = aInfo.height
		}
		connRow := topRow + maxH

		// Draw connector bar linking function output to input output.
		g.hLine(connRow, funOutCol, argOutCol)

		// Draw vertical connectors from sub-term outputs down to connector bar.
		if funOutRow < connRow {
			g.vLine(funOutCol, funOutRow, connRow)
		}
		if argOutRow < connRow {
			g.vLine(argOutCol, argOutRow, connRow)
		}

		// Output of application = function's root column, at connector row.
		outRow = connRow
		outCol = funOutCol
		return
	}
	return topRow, leftCol
}
