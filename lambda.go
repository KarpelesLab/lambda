package lambda

import (
	"fmt"
)

// Object is the interface for all lambda calculus terms
type Object interface {
	String() string
	// FreeVars returns the set of free variables in the term
	FreeVars() map[string]bool
	// Substitute replaces a variable with a term
	Substitute(varName string, replacement Object) Object
	// AlphaConvert renames a bound variable
	AlphaConvert(oldName, newName string) Object
	// BetaReduce performs one step of β-reduction if possible
	BetaReduce() (Object, bool)
	// EtaConvert performs η-conversion if possible
	EtaConvert() (Object, bool)
}

// Var represents a variable
type Var struct {
	Name string
}

// Abstraction represents an abstraction (λx.t)
type Abstraction struct {
	Param string // The bound variable
	Body  Object // The body of the abstraction
}

// Application represents an application (t s)
type Application struct {
	Func Object // The function
	Arg  Object // The argument
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
func (v Var) Substitute(varName string, replacement Object) Object {
	if v.Name == varName {
		return replacement
	}
	return v
}

func (a Abstraction) Substitute(varName string, replacement Object) Object {
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

func (a Application) Substitute(varName string, replacement Object) Object {
	return Application{
		Func: a.Func.Substitute(varName, replacement),
		Arg:  a.Arg.Substitute(varName, replacement),
	}
}

// AlphaConvert implementations
func (v Var) AlphaConvert(oldName, newName string) Object {
	if v.Name == oldName {
		return Var{Name: newName}
	}
	return v
}

func (a Abstraction) AlphaConvert(oldName, newName string) Object {
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

func (a Application) AlphaConvert(oldName, newName string) Object {
	return Application{
		Func: a.Func.AlphaConvert(oldName, newName),
		Arg:  a.Arg.AlphaConvert(oldName, newName),
	}
}

// Reduce performs multiple β-reductions up to a maximum number of steps.
// It returns the reduced term and the number of reductions performed.
// If limit is 0 or negative, a default limit of 1000 is used.
func Reduce(obj Object, limit int) (Object, int) {
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
func (v Var) BetaReduce() (Object, bool) {
	return v, false
}

func (a Abstraction) BetaReduce() (Object, bool) {
	// Try to reduce the body
	newBody, reduced := a.Body.BetaReduce()
	if reduced {
		return Abstraction{Param: a.Param, Body: newBody}, true
	}
	return a, false
}

func (a Application) BetaReduce() (Object, bool) {
	// Check if we can do β-reduction at the top level
	if abs, ok := a.Func.(Abstraction); ok {
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
func (v Var) EtaConvert() (Object, bool) {
	return v, false
}

func (a Abstraction) EtaConvert() (Object, bool) {
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

func (a Application) EtaConvert() (Object, bool) {
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
func ChurchNumeral(n int) Object {
	if n < 0 {
		panic("Church numerals are only defined for non-negative integers")
	}

	// Build f (f (f ... (f x)...))
	var body Object = Var{Name: "x"}
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
func ToInt(term Object) int {
	// Apply the Church numeral to an increment function and 0
	// Church numeral n applied to f and 0 will call f n times
	// We'll use a marker to count
	result := Object(Application{
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
func ToBool(term Object) bool {
	// Apply the boolean to two distinct markers
	// TRUE will return the first argument, FALSE will return the second
	result := Object(Application{
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
func countApplications(term Object, funcName string) int {
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