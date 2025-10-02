package lambda

import (
	"fmt"
)

// Term is the interface for all lambda calculus terms
type Term interface {
	String() string
	// FreeVars returns the set of free variables in the term
	FreeVars() map[string]bool
	// Substitute replaces a variable with a term
	Substitute(varName string, replacement Term) Term
	// AlphaConvert renames a bound variable
	AlphaConvert(oldName, newName string) Term
	// BetaReduce performs one step of β-reduction if possible
	BetaReduce() (Term, bool)
	// EtaConvert performs η-conversion if possible
	EtaConvert() (Term, bool)
}

// LazyScript holds an unparsed expression that will be parsed on first use.
// This allows defining constants that use Parse without creating initialization cycles.
type LazyScript struct {
	script string
	parsed Term
}

// MakeLazyScript creates a new LazyScript from a string
func MakeLazyScript(script string) *LazyScript {
	return &LazyScript{script: script}
}

// parse parses and caches the expression on first use
func (l *LazyScript) parse() Term {
	if l.parsed == nil {
		parsed, err := Parse(l.script)
		if err != nil {
			panic(fmt.Sprintf("LazyScript parse error: %v\nScript: %s", err, l.script))
		}
		l.parsed = parsed
	}
	return l.parsed
}

func (l *LazyScript) String() string {
	return l.parse().String()
}

func (l *LazyScript) FreeVars() map[string]bool {
	return l.parse().FreeVars()
}

func (l *LazyScript) Substitute(varName string, replacement Term) Term {
	return l.parse().Substitute(varName, replacement)
}

func (l *LazyScript) AlphaConvert(oldName, newName string) Term {
	return l.parse().AlphaConvert(oldName, newName)
}

func (l *LazyScript) BetaReduce() (Term, bool) {
	return l.parse().BetaReduce()
}

func (l *LazyScript) EtaConvert() (Term, bool) {
	return l.parse().EtaConvert()
}

// Var represents a variable
type Var struct {
	Name string
}

// Abstraction represents an abstraction (λx.t)
type Abstraction struct {
	Param string // The bound variable
	Body  Term   // The body of the abstraction
}

// Application represents an application (t s)
type Application struct {
	Func Term // The function
	Arg  Term // The argument
}

// String methods
func (v Var) String() string {
	return v.Name
}

func (a Abstraction) String() string {
	return fmt.Sprintf("λ%s.%s", a.Param, a.Body.String())
}

func (a Application) String() string {
	// Add parentheses when necessary
	funcStr := a.Func.String()
	if _, isAbs := a.Func.(Abstraction); isAbs {
		funcStr = "(" + funcStr + ")"
	}

	argStr := a.Arg.String()
	if _, isApp := a.Arg.(Application); isApp {
		argStr = "(" + argStr + ")"
	} else if _, isAbs := a.Arg.(Abstraction); isAbs {
		argStr = "(" + argStr + ")"
	}

	return funcStr + " " + argStr
}

// FreeVars implementations
func (v Var) FreeVars() map[string]bool {
	return map[string]bool{v.Name: true}
}

func (a Abstraction) FreeVars() map[string]bool {
	fv := a.Body.FreeVars()
	delete(fv, a.Param)
	return fv
}

func (a Application) FreeVars() map[string]bool {
	fv := a.Func.FreeVars()
	for k := range a.Arg.FreeVars() {
		fv[k] = true
	}
	return fv
}

// Substitute implementations
func (v Var) Substitute(varName string, replacement Term) Term {
	if v.Name == varName {
		return replacement
	}
	return v
}

func (a Abstraction) Substitute(varName string, replacement Term) Term {
	if a.Param == varName {
		// Variable is bound, no substitution in body
		return a
	}

	// Check for variable capture
	if replacement.FreeVars()[a.Param] {
		// Need α-conversion to avoid capture
		newParam := freshVar(a.Param, replacement.FreeVars())
		newBody := a.Body.AlphaConvert(a.Param, newParam)
		return Abstraction{Param: newParam, Body: newBody.Substitute(varName, replacement)}
	}

	return Abstraction{Param: a.Param, Body: a.Body.Substitute(varName, replacement)}
}

func (a Application) Substitute(varName string, replacement Term) Term {
	return Application{
		Func: a.Func.Substitute(varName, replacement),
		Arg:  a.Arg.Substitute(varName, replacement),
	}
}

// AlphaConvert implementations
func (v Var) AlphaConvert(oldName, newName string) Term {
	if v.Name == oldName {
		return Var{Name: newName}
	}
	return v
}

func (a Abstraction) AlphaConvert(oldName, newName string) Term {
	if a.Param == oldName {
		return Abstraction{
			Param: newName,
			Body:  a.Body.AlphaConvert(oldName, newName),
		}
	}
	return Abstraction{
		Param: a.Param,
		Body:  a.Body.AlphaConvert(oldName, newName),
	}
}

func (a Application) AlphaConvert(oldName, newName string) Term {
	return Application{
		Func: a.Func.AlphaConvert(oldName, newName),
		Arg:  a.Arg.AlphaConvert(oldName, newName),
	}
}

// Reduce performs multiple β-reductions up to a maximum number of steps.
// It returns the reduced term and the number of reductions performed.
// If limit is 0 or negative, a default limit of 1000 is used.
func Reduce(obj Term, limit int) (Term, int) {
	if limit <= 0 {
		limit = 1000
	}

	steps := 0
	for i := 0; i < limit; i++ {
		reduced, didReduce := obj.BetaReduce()
		if !didReduce {
			break
		}
		obj = reduced
		steps++
	}

	return obj, steps
}

// BetaReduce implementations
func (v Var) BetaReduce() (Term, bool) {
	return v, false
}

func (a Abstraction) BetaReduce() (Term, bool) {
	// Try to reduce the body
	newBody, reduced := a.Body.BetaReduce()
	if reduced {
		return Abstraction{Param: a.Param, Body: newBody}, true
	}
	return a, false
}

func (a Application) BetaReduce() (Term, bool) {
	// Unwrap LazyScript if present
	funcTerm := a.Func
	if ls, ok := funcTerm.(*LazyScript); ok {
		funcTerm = ls.parse()
	}

	// Check if we can do β-reduction at the top level
	if abs, ok := funcTerm.(Abstraction); ok {
		// (λx.t) s → t[x := s]
		result := abs.Body.Substitute(abs.Param, a.Arg)
		return result, true
	}

	// Try to reduce the function
	newFunc, reduced := a.Func.BetaReduce()
	if reduced {
		return Application{Func: newFunc, Arg: a.Arg}, true
	}

	// Try to reduce the argument
	newArg, reduced := a.Arg.BetaReduce()
	if reduced {
		return Application{Func: a.Func, Arg: newArg}, true
	}

	return a, false
}

// EtaConvert implementations
func (v Var) EtaConvert() (Term, bool) {
	return v, false
}

func (a Abstraction) EtaConvert() (Term, bool) {
	// η-reduction: λx.(f x) → f if x is not free in f
	if app, ok := a.Body.(Application); ok {
		if v, ok := app.Arg.(Var); ok && v.Name == a.Param {
			// Check if a.Param is not free in app.Func
			if !app.Func.FreeVars()[a.Param] {
				return app.Func, true
			}
		}
	}

	// Try to η-convert the body
	newBody, converted := a.Body.EtaConvert()
	if converted {
		return Abstraction{Param: a.Param, Body: newBody}, true
	}

	return a, false
}

func (a Application) EtaConvert() (Term, bool) {
	// Try to η-convert the function
	newFunc, converted := a.Func.EtaConvert()
	if converted {
		return Application{Func: newFunc, Arg: a.Arg}, true
	}

	// Try to η-convert the argument
	newArg, converted := a.Arg.EtaConvert()
	if converted {
		return Application{Func: a.Func, Arg: newArg}, true
	}

	return a, false
}

// Helper function to generate fresh variable names
func freshVar(base string, avoid map[string]bool) string {
	if !avoid[base] {
		return base
	}
	i := 0
	for {
		candidate := fmt.Sprintf("%s%d", base, i)
		if !avoid[candidate] {
			return candidate
		}
		i++
	}
}

// ChurchNumeral creates a Church numeral for the given natural number
// n := λf.λx.f^n x
// 0 := λf.λx.x
// 1 := λf.λx.f x
// 2 := λf.λx.f (f x)
func ChurchNumeral(n int) Term {
	if n < 0 {
		panic("Church numerals are only defined for non-negative integers")
	}

	// Build f (f (f ... (f x)...))
	var body Term = Var{Name: "x"}
	for i := 0; i < n; i++ {
		body = Application{Func: Var{Name: "f"}, Arg: body}
	}

	// Wrap in λx and λf
	return Abstraction{
		Param: "f",
		Body: Abstraction{
			Param: "x",
			Body:  body,
		},
	}
}

// ToInt converts a Church numeral to a Go integer by applying it to increment and 0
// Church numeral n = λf.λx.f^n x, so we apply it to a marker function and count applications
func ToInt(term Term) int {
	// Apply the Church numeral to an increment function and 0
	// Church numeral n applied to f and 0 will call f n times
	// We'll use a marker to count
	result := Term(Application{
		Func: Application{
			Func: term,
			Arg:  Var{Name: "SUCC_MARKER"},
		},
		Arg: Var{Name: "ZERO_MARKER"},
	})

	// Reduce completely (with a limit to avoid infinite loops)
	result, _ = Reduce(result, 1000)

	// Count nested applications of SUCC_MARKER
	return countApplications(result, "SUCC_MARKER")
}

// ToBool converts a Church boolean to a Go bool
// TRUE returns true, FALSE returns false
// Church boolean: TRUE = λx.λy.x, FALSE = λx.λy.y
func ToBool(term Term) bool {
	// Apply the boolean to two distinct markers
	// TRUE will return the first argument, FALSE will return the second
	result := Term(Application{
		Func: Application{
			Func: term,
			Arg:  Var{Name: "TRUE_MARKER"},
		},
		Arg: Var{Name: "FALSE_MARKER"},
	})

	// Reduce completely (with a limit to avoid infinite loops)
	result, _ = Reduce(result, 1000)

	// Check which marker we got
	if v, ok := result.(Var); ok {
		if v.Name == "TRUE_MARKER" {
			return true
		}
		if v.Name == "FALSE_MARKER" {
			return false
		}
	}

	// Default to false if we can't determine
	return false
}

// Helper function to count nested applications of a specific function
func countApplications(term Term, funcName string) int {
	switch t := term.(type) {
	case Var:
		if t.Name == "ZERO_MARKER" {
			return 0
		}
		return 0
	case Application:
		if v, ok := t.Func.(Var); ok && v.Name == funcName {
			return 1 + countApplications(t.Arg, funcName)
		}
		// Otherwise, try to count in the argument
		return countApplications(t.Arg, funcName)
	case Abstraction:
		return countApplications(t.Body, funcName)
	}
	return 0
}