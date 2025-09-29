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

// Standard combinators
//
// I is the identity function.
//
// SK and BCKW form complete combinator calculus systems that can express any lambda term.
// This means that any lambda calculus expression can be translated into an equivalent expression
// using only these combinators.
//
// Ω is UU (or ω ω), the smallest term that has no normal form - it reduces to itself infinitely.
// YI is another such term with no normal form.
var (
	// I := λx.x (Identity function)
	I = Abstraction{
		Param: "x",
		Body:  Var{Name: "x"},
	}

	// K := λx.λy.x (Constant/Cancel)
	// Together with S, forms a complete combinator calculus basis (SK calculus)
	K = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body:  Var{Name: "x"},
		},
	}

	// S := λx.λy.λz.x z (y z) (Substitution)
	// Together with K, forms a complete combinator calculus basis (SK calculus)
	S = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body: Abstraction{
				Param: "z",
				Body: Application{
					Func: Application{
						Func: Var{Name: "x"},
						Arg:  Var{Name: "z"},
					},
					Arg: Application{
						Func: Var{Name: "y"},
						Arg:  Var{Name: "z"},
					},
				},
			},
		},
	}

	// B := λx.λy.λz.x (y z) (Composition)
	// Together with C, K, and W, forms a complete combinator calculus basis (BCKW calculus)
	B = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body: Abstraction{
				Param: "z",
				Body: Application{
					Func: Var{Name: "x"},
					Arg: Application{
						Func: Var{Name: "y"},
						Arg:  Var{Name: "z"},
					},
				},
			},
		},
	}

	// C := λx.λy.λz.x z y (Flip)
	// Together with B, K, and W, forms a complete combinator calculus basis (BCKW calculus)
	C = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body: Abstraction{
				Param: "z",
				Body: Application{
					Func: Application{
						Func: Var{Name: "x"},
						Arg:  Var{Name: "z"},
					},
					Arg: Var{Name: "y"},
				},
			},
		},
	}

	// W := λx.λy.x y y (Warbler/Duplication)
	// Together with B, C, and K, forms a complete combinator calculus basis (BCKW calculus)
	W = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body: Application{
				Func: Application{
					Func: Var{Name: "x"},
					Arg:  Var{Name: "y"},
				},
				Arg: Var{Name: "y"},
			},
		},
	}

	// U := λx.x x (Self-application)
	// Also known as ω (omega) or Δ (delta)
	U = Abstraction{
		Param: "x",
		Body: Application{
			Func: Var{Name: "x"},
			Arg:  Var{Name: "x"},
		},
	}

	// Ω (Omega) := U U (or ω ω)
	// The smallest term that has no normal form - it reduces to itself infinitely
	// Another example of a term with no normal form is Y I
	OMEGA = Application{
		Func: U,
		Arg:  U,
	}
)

// Aliases for combinators
var (
	OMEGA_LOWER = U     // ω := λx.x x (same as U)
	DELTA       = U     // δ := λx.x x (same as U)
)

// Boolean constants
//
// TRUE and FALSE are commonly abbreviated as T and F.
var (
	// TRUE := λx.λy.x (same as K combinator)
	// Commonly abbreviated as T
	TRUE = K

	// FALSE := λx.λy.y
	// Commonly abbreviated as F
	FALSE = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body:  Var{Name: "y"},
		},
	}

	// T is an alias for TRUE
	T = TRUE

	// F is an alias for FALSE
	F = FALSE
)

// Boolean operations
var (
	// AND := λp.λq.p q p
	AND = Abstraction{
		Param: "p",
		Body: Abstraction{
			Param: "q",
			Body: Application{
				Func: Application{
					Func: Var{Name: "p"},
					Arg:  Var{Name: "q"},
				},
				Arg: Var{Name: "p"},
			},
		},
	}

	// OR := λp.λq.p p q
	OR = Abstraction{
		Param: "p",
		Body: Abstraction{
			Param: "q",
			Body: Application{
				Func: Application{
					Func: Var{Name: "p"},
					Arg:  Var{Name: "p"},
				},
				Arg: Var{Name: "q"},
			},
		},
	}

	// NOT := λp.p FALSE TRUE
	NOT = Abstraction{
		Param: "p",
		Body: Application{
			Func: Application{
				Func: Var{Name: "p"},
				Arg:  FALSE,
			},
			Arg: TRUE,
		},
	}
)

// Control flow
var (
	// IFTHENELSE := λp.λa.λb.p a b
	IFTHENELSE = Abstraction{
		Param: "p",
		Body: Abstraction{
			Param: "a",
			Body: Abstraction{
				Param: "b",
				Body: Application{
					Func: Application{
						Func: Var{Name: "p"},
						Arg:  Var{Name: "a"},
					},
					Arg: Var{Name: "b"},
				},
			},
		},
	}
)

// Arithmetic operations
var (
	// SUCC := λn.λf.λx.f (n f x)
	SUCC = Abstraction{
		Param: "n",
		Body: Abstraction{
			Param: "f",
			Body: Abstraction{
				Param: "x",
				Body: Application{
					Func: Var{Name: "f"},
					Arg: Application{
						Func: Application{
							Func: Var{Name: "n"},
							Arg:  Var{Name: "f"},
						},
						Arg: Var{Name: "x"},
					},
				},
			},
		},
	}

	// PLUS := λm.λn.λf.λx.m f (n f x)
	PLUS = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Abstraction{
				Param: "f",
				Body: Abstraction{
					Param: "x",
					Body: Application{
						Func: Application{
							Func: Var{Name: "m"},
							Arg:  Var{Name: "f"},
						},
						Arg: Application{
							Func: Application{
								Func: Var{Name: "n"},
								Arg:  Var{Name: "f"},
							},
							Arg: Var{Name: "x"},
						},
					},
				},
			},
		},
	}

	// SUB := λm.λn.n PRED m
	SUB = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg:  Var{Name: "PRED"},
				},
				Arg: Var{Name: "m"},
			},
		},
	}

	// MULT := λm.λn.λf.m (n f)
	MULT = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Abstraction{
				Param: "f",
				Body: Application{
					Func: Var{Name: "m"},
					Arg: Application{
						Func: Var{Name: "n"},
						Arg:  Var{Name: "f"},
					},
				},
			},
		},
	}
)

// Predicates
var (
	// ISZERO := λn.n (λx.FALSE) TRUE
	ISZERO = Abstraction{
		Param: "n",
		Body: Application{
			Func: Application{
				Func: Var{Name: "n"},
				Arg: Abstraction{
					Param: "x",
					Body:  FALSE,
				},
			},
			Arg: TRUE,
		},
	}

	// LEQ := λm.λn.ISZERO (SUB m n)
	LEQ = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: ISZERO,
				Arg: Application{
					Func: Application{
						Func: SUB,
						Arg:  Var{Name: "m"},
					},
					Arg: Var{Name: "n"},
				},
			},
		},
	}
)

// Pair operations
var (
	// PAIR := λx.λy.λf.f x y
	PAIR = Abstraction{
		Param: "x",
		Body: Abstraction{
			Param: "y",
			Body: Abstraction{
				Param: "f",
				Body: Application{
					Func: Application{
						Func: Var{Name: "f"},
						Arg:  Var{Name: "x"},
					},
					Arg: Var{Name: "y"},
				},
			},
		},
	}

	// FIRST := λp.p TRUE
	FIRST = Abstraction{
		Param: "p",
		Body: Application{
			Func: Var{Name: "p"},
			Arg:  TRUE,
		},
	}

	// SECOND := λp.p FALSE
	SECOND = Abstraction{
		Param: "p",
		Body: Application{
			Func: Var{Name: "p"},
			Arg:  FALSE,
		},
	}
)

// Φ combinator for PRED
var (
	// Φ := λx.PAIR (SECOND x) (SUCC (SECOND x))
	PHI = Abstraction{
		Param: "x",
		Body: Application{
			Func: Application{
				Func: PAIR,
				Arg: Application{
					Func: SECOND,
					Arg:  Var{Name: "x"},
				},
			},
			Arg: Application{
				Func: SUCC,
				Arg: Application{
					Func: SECOND,
					Arg:  Var{Name: "x"},
				},
			},
		},
	}

	// PRED := λn.FIRST (n Φ (PAIR 0 0))
	PRED = Abstraction{
		Param: "n",
		Body: Application{
			Func: FIRST,
			Arg: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg:  PHI,
				},
				Arg: Application{
					Func: Application{
						Func: PAIR,
						Arg:  ChurchNumeral(0),
					},
					Arg: ChurchNumeral(0),
				},
			},
		},
	}
)

// List operations
var (
	// NIL := λx.TRUE
	NIL = Abstraction{
		Param: "x",
		Body:  TRUE,
	}

	// NULL := λp.p (λx.λy.FALSE)
	NULL = Abstraction{
		Param: "p",
		Body: Application{
			Func: Var{Name: "p"},
			Arg: Abstraction{
				Param: "x",
				Body: Abstraction{
					Param: "y",
					Body:  FALSE,
				},
			},
		},
	}
)

// Y combinator for recursion
//
// Y := λf.(λx.f (x x)) (λx.f (x x))
//
// The Y combinator enables recursion in lambda calculus.
// It satisfies the property: Y g = g (Y g)
//
// Alternative definition: Y = B U (C B U)
// This alternative shows Y in terms of B, C, and U combinators.
var Y = Abstraction{
	Param: "f",
	Body: Application{
		Func: Abstraction{
			Param: "x",
			Body: Application{
				Func: Var{Name: "f"},
				Arg: Application{
					Func: Var{Name: "x"},
					Arg:  Var{Name: "x"},
				},
			},
		},
		Arg: Abstraction{
			Param: "x",
			Body: Application{
				Func: Var{Name: "f"},
				Arg: Application{
					Func: Var{Name: "x"},
					Arg:  Var{Name: "x"},
				},
			},
		},
	},
}

// FACTORIAL := Y (λf.λn.ISZERO n 1 (MULT n (f (PRED n))))
var FACTORIAL = Application{
	Func: Y,
	Arg: Abstraction{
		Param: "f",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: Application{
					Func: Application{
						Func: ISZERO,
						Arg:  Var{Name: "n"},
					},
					Arg: ChurchNumeral(1),
				},
				Arg: Application{
					Func: Application{
						Func: MULT,
						Arg:  Var{Name: "n"},
					},
					Arg: Application{
						Func: Var{Name: "f"},
						Arg: Application{
							Func: PRED,
							Arg:  Var{Name: "n"},
						},
					},
				},
			},
		},
	},
}

// FAC is an alternative factorial implementation without Y combinator
// FAC = λn.λf.n(λf.λn.n(f(λf.λx.n f(f x))))(λx.f)(λx.x)
var FAC = Abstraction{
	Param: "n",
	Body: Abstraction{
		Param: "f",
		Body: Application{
			Func: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg: Abstraction{
						Param: "f",
						Body: Abstraction{
							Param: "n",
							Body: Application{
								Func: Var{Name: "n"},
								Arg: Application{
									Func: Var{Name: "f"},
									Arg: Abstraction{
										Param: "f",
										Body: Abstraction{
											Param: "x",
											Body: Application{
												Func: Application{
													Func: Var{Name: "n"},
													Arg:  Var{Name: "f"},
												},
												Arg: Application{
													Func: Var{Name: "f"},
													Arg:  Var{Name: "x"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Arg: Abstraction{
					Param: "x",
					Body:  Var{Name: "f"},
				},
			},
			Arg: Abstraction{
				Param: "x",
				Body:  Var{Name: "x"},
			},
		},
	},
}

// FIB is a Fibonacci implementation without Y combinator
// FIB = λn.λf.n(λc.λa.λb.c b(λx.a (b x)))(λx.λy.x)(λx.x)f
var FIB = Abstraction{
	Param: "n",
	Body: Abstraction{
		Param: "f",
		Body: Application{
			Func: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg: Abstraction{
						Param: "c",
						Body: Abstraction{
							Param: "a",
							Body: Abstraction{
								Param: "b",
								Body: Application{
									Func: Application{
										Func: Var{Name: "c"},
										Arg:  Var{Name: "b"},
									},
									Arg: Abstraction{
										Param: "x",
										Body: Application{
											Func: Var{Name: "a"},
											Arg: Application{
												Func: Var{Name: "b"},
												Arg:  Var{Name: "x"},
											},
										},
									},
								},
							},
						},
					},
				},
				Arg: Abstraction{
					Param: "x",
					Body: Abstraction{
						Param: "y",
						Body: Var{Name: "x"},
					},
				},
			},
			Arg: Abstraction{
				Param: "x",
				Body: Var{Name: "x"},
			},
		},
	},
}

// ToInt converts a Church numeral to a Go integer by applying it to increment and 0
func ToInt(term Object) int {
	// Reduce the term as much as possible first
	for {
		reduced, didReduce := term.BetaReduce()
		if !didReduce {
			break
		}
		term = reduced
	}

	// Apply the Church numeral to an increment function and 0
	// We'll track how many times the function is called
	count := 0

	// Create a simple evaluator that counts applications
	// Church numeral n applied to f and 0 will call f n times
	// We'll use a marker to count

	// Apply to a successor-like function: λx.x+1
	// We'll represent numbers as nested applications of a marker
	var result Object = Application{
		Func: Application{
			Func: term,
			Arg:  Var{Name: "SUCC_MARKER"},
		},
		Arg: Var{Name: "ZERO_MARKER"},
	}

	// Reduce completely
	for i := 0; i < 1000; i++ {
		reduced, didReduce := result.BetaReduce()
		if !didReduce {
			break
		}
		result = reduced
	}

	// Count nested applications of SUCC_MARKER
	count = countApplications(result, "SUCC_MARKER")

	return count
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