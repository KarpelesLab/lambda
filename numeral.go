package lambda

import "fmt"

// Numeral is a compact Church numeral representation backed by uint64.
// It behaves identically to λf.λx.f^n(x) during beta reduction but
// avoids the O(n) AST. When applied to an argument f, it produces a
// NumeralApply (partially applied). When that is applied to x, it
// expands f^n(x) in a single reduction step.
type Numeral uint64

func (n Numeral) String() string {
	return fmt.Sprintf("[%d]", uint64(n))
}

func (n Numeral) FreeVars() map[string]bool {
	return make(map[string]bool)
}

func (n Numeral) Substitute(string, Term) Term {
	return n
}

func (n Numeral) AlphaConvert(string, string) Term {
	return n
}

func (n Numeral) BetaReduce() (Term, bool) {
	return n, false
}

func (n Numeral) EtaConvert() (Term, bool) {
	return n, false
}

// Expand converts a Numeral to its full Church encoding λf.λx.f^n(x).
func (n Numeral) Expand() Term {
	var body Term = Var{Name: "x"}
	for i := uint64(0); i < uint64(n); i++ {
		body = Application{Func: Var{Name: "f"}, Arg: body}
	}
	return Abstraction{
		Param: "f",
		Body:  Abstraction{Param: "x", Body: body},
	}
}

// NumeralApply is a partially applied Church numeral: λparam.F^N(param).
// Created when a Numeral is applied to a function argument.
// When applied to x, it expands F^N(x) in a single reduction step.
type NumeralApply struct {
	N     uint64
	Param string
	F     Term
}

func (na NumeralApply) String() string {
	if na.N == 0 {
		return fmt.Sprintf("λ%s.%s", na.Param, na.Param)
	}
	if na.N == 1 {
		return fmt.Sprintf("λ%s.%s %s", na.Param, na.F.String(), na.Param)
	}
	return fmt.Sprintf("λ%s.%s^%d %s", na.Param, na.F.String(), na.N, na.Param)
}

func (na NumeralApply) FreeVars() map[string]bool {
	fv := na.F.FreeVars()
	// Param is bound — by construction F does not contain Param free,
	// but delete it defensively.
	delete(fv, na.Param)
	return fv
}

func (na NumeralApply) Substitute(varName string, replacement Term) Term {
	if varName == na.Param {
		return na
	}
	if replacement.FreeVars()[na.Param] {
		avoid := make(map[string]bool)
		for k := range replacement.FreeVars() {
			avoid[k] = true
		}
		for k := range na.F.FreeVars() {
			avoid[k] = true
		}
		newParam := freshVar(na.Param, avoid)
		// F doesn't contain Param as free, so no alpha-rename needed in F
		return NumeralApply{
			N:     na.N,
			Param: newParam,
			F:     na.F.Substitute(varName, replacement),
		}
	}
	return NumeralApply{N: na.N, Param: na.Param, F: na.F.Substitute(varName, replacement)}
}

func (na NumeralApply) AlphaConvert(oldName, newName string) Term {
	if na.Param == oldName {
		return NumeralApply{N: na.N, Param: newName, F: na.F}
	}
	return NumeralApply{N: na.N, Param: na.Param, F: na.F.AlphaConvert(oldName, newName)}
}

func (na NumeralApply) BetaReduce() (Term, bool) {
	newF, reduced := na.F.BetaReduce()
	if reduced {
		return NumeralApply{N: na.N, Param: na.Param, F: newF}, true
	}
	return na, false
}

func (na NumeralApply) EtaConvert() (Term, bool) {
	// λx.f^1(x) = λx.f x → f  (η-reduction, if x not free in f)
	if na.N == 1 && !na.F.FreeVars()[na.Param] {
		return na.F, true
	}
	return na, false
}

// expand builds the full application chain F^n(x), collapsing nested
// NumeralApplys first (multiplication optimization).
func (na NumeralApply) expand(x Term) Term {
	f := na.F
	n := na.N

	// Flatten nested NumeralApplys:
	// NumeralApply{a, NumeralApply{b, g}}.expand(x) = g^(a*b)(x)
	// because (g^b)^a = g^(a*b)
	for {
		if inner, ok := f.(NumeralApply); ok {
			n *= inner.N
			f = inner.F
		} else {
			break
		}
	}

	result := x
	for i := uint64(0); i < n; i++ {
		result = Application{Func: f, Arg: result}
	}
	return result
}
