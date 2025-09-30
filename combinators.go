package lambda

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
	I = MakeLazyScript(`λx.x`)

	// K := λx.λy.x (Constant/Cancel)
	// Together with S, forms a complete combinator calculus basis (SK calculus)
	K = MakeLazyScript(`λx.λy.x`)

	// S := λx.λy.λz.x z (y z) (Substitution)
	// Together with K, forms a complete combinator calculus basis (SK calculus)
	S = MakeLazyScript(`λx.λy.λz.x z (y z)`)

	// B := λx.λy.λz.x (y z) (Composition)
	// Together with C, K, and W, forms a complete combinator calculus basis (BCKW calculus)
	B = MakeLazyScript(`λx.λy.λz.x (y z)`)

	// C := λx.λy.λz.x z y (Flip)
	// Together with B, K, and W, forms a complete combinator calculus basis (BCKW calculus)
	C = MakeLazyScript(`λx.λy.λz.x z y`)

	// W := λx.λy.x y y (Warbler/Duplication)
	// Together with B, C, and K, forms a complete combinator calculus basis (BCKW calculus)
	W = MakeLazyScript(`λx.λy.x y y`)

	// U := λx.x x (Self-application)
	// Also known as ω (omega) or Δ (delta)
	U = MakeLazyScript(`λx.x x`)

	// Ω (Omega) := U U (or ω ω)
	// The smallest term that has no normal form - it reduces to itself infinitely
	// Another example of a term with no normal form is Y I
	OMEGA = MakeLazyScript(`(λx.x x) (λx.x x)`)
)

// Aliases for combinators
var (
	OMEGA_LOWER = U // ω := λx.x x (same as U)
	DELTA       = U // δ := λx.x x (same as U)
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
	FALSE = MakeLazyScript(`λx.λy.y`)

	// T is an alias for TRUE
	T = TRUE

	// F is an alias for FALSE
	F = FALSE
)

// Boolean operations
var (
	// AND := λp.λq.p q p
	AND = MakeLazyScript(`λp.λq.p q p`)

	// OR := λp.λq.p p q
	OR = MakeLazyScript(`λp.λq.p p q`)

	// NOT := λp.p FALSE TRUE
	NOT = MakeLazyScript(`λp.p _FALSE _TRUE`)
)

// Control flow
var (
	// IF := λb.λx.λy.b x y
	IF = MakeLazyScript(`λb.λx.λy.b x y`)

	// IFTHENELSE := λp.λa.λb.p a b (same as IF)
	IFTHENELSE = IF
)

// Arithmetic operations
var (
	// ZERO := λf.λx.x (Church numeral 0)
	ZERO = MakeLazyScript(`λf.λx.x`)

	// SUCC := λn.λf.λx.f (n f x)
	SUCC = MakeLazyScript(`λn.λf.λx.f (n f x)`)

	// ONE := SUCC ZERO (Church numeral 1)
	ONE = MakeLazyScript(`_SUCC _ZERO`)

	// TWO := SUCC ONE (Church numeral 2)
	TWO = MakeLazyScript(`_SUCC _ONE`)

	// THREE := SUCC TWO (Church numeral 3)
	THREE = MakeLazyScript(`_SUCC _TWO`)

	// DEC := PRED (decrement, alias for predecessor)
	DEC = PRED

	// ADD := PLUS (addition, alias)
	ADD = PLUS

	// PLUS := λm.λn.λf.λx.m f (n f x)
	PLUS = MakeLazyScript(`λm.λn.λf.λx.m f (n f x)`)

	// SUB := λm.λn.n PRED m
	SUB = MakeLazyScript(`λm.λn.n _PRED m`)

	// MULT := λm.λn.λf.m (n f)
	MULT = MakeLazyScript(`λm.λn.λf.m (n f)`)

	// POW := λb.λn.n b (exponentiation: b^n)
	POW = MakeLazyScript(`λb.λn.n b`)
)

// Predicates
var (
	// ISZERO := λn.n (λx.FALSE) TRUE
	ISZERO = MakeLazyScript(`λn.n (λx._FALSE) _TRUE`)

	// LEQ := λm.λn.ISZERO (SUB m n)
	LEQ = MakeLazyScript(`λm.λn._ISZERO (_SUB m n)`)

	// LT := λm.λn.NOT (LEQ n m)
	LT = MakeLazyScript(`λm.λn._NOT (_LEQ n m)`)

	// EQ := λm.λn.AND (LEQ m n) (LEQ n m)
	EQ = MakeLazyScript(`λm.λn._AND (_LEQ m n) (_LEQ n m)`)

	// MAX := λa.λb.IF (LEQ a b) b a
	MAX = MakeLazyScript(`λa.λb._IF (_LEQ a b) b a`)

	// MIN := λa.λb.IF (LEQ a b) a b
	MIN = MakeLazyScript(`λa.λb._IF (_LEQ a b) a b`)
)

// GCD := Y (λrec.λa.λb.IF (ISZERO b) a (rec b (MOD a b)))
var GCD = MakeLazyScript(`
	_Y (λrec.λa.λb.
		_IF (_ISZERO b)
			a
			(rec b (_MOD a b)))
`)

// MOD := Y (λrec.λm.λn.(ISZERO n) ZERO ((LT m n) m (rec (SUB m n) n)))
// Modulo operation with zero-divisor guard: m mod n = 0 if n = 0
var MOD = MakeLazyScript(`
	_Y (λrec.λm.λn.
		(_ISZERO n) _ZERO
		((_LT m n) m (rec (_SUB m n) n)))
`)

// Pair operations
var (
	// PAIR := λx.λy.λf.f x y
	PAIR = MakeLazyScript(`λx.λy.λf.f x y`)

	// FIRST := λp.p TRUE
	FIRST = MakeLazyScript(`λp.p _TRUE`)

	// SECOND := λp.p FALSE
	SECOND = MakeLazyScript(`λp.p _FALSE`)
)

// Bit manipulation helpers
var (
	// STEP2 := λp.PAIR (IF (SECOND p) (SUCC (FIRST p)) (FIRST p)) (NOT (SECOND p))
	STEP2 = MakeLazyScript(`λp._PAIR (_IF (_SECOND p) (_SUCC (_FIRST p)) (_FIRST p)) (_NOT (_SECOND p))`)

	// INIT2 := PAIR ZERO FALSE
	INIT2 = MakeLazyScript(`_PAIR _ZERO _FALSE`)
)

// Division and parity operations
var (
	// DIV2 := λn.FIRST (n STEP2 INIT2)
	DIV2 = MakeLazyScript(`λn._FIRST (n _STEP2 _INIT2)`)

	// ISODD := λn.SECOND (n STEP2 INIT2)
	ISODD = MakeLazyScript(`λn._SECOND (n _STEP2 _INIT2)`)

	// ISEVEN := λn.NOT (ISODD n)
	ISEVEN = MakeLazyScript(`λn._NOT (_ISODD n)`)
)

// MUL := λm.λn.λf.m (n f)
// Note: MUL is already defined above as MULT, but we need it for POWMOD
var MUL = MULT

// POWMOD' := Y (λrec.λa.λe.λm.λr.IF (ISZERO e) (IF (ISZERO m) r (MOD r m)) (IF (ISEVEN e) (rec (MOD (MUL a a) m) (DIV2 e) m r) (rec (MOD (MUL a a) m) (DIV2 e) m (MOD (MUL r a) m))))
// Tail-recursive modular exponentiation with accumulator
var POWMOD_PRIME = MakeLazyScript(`
	_Y (λrec.λa.λe.λm.λr.
		_IF (_ISZERO e)
			(_IF (_ISZERO m) r (_MOD r m))
			(_IF (_ISEVEN e)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m r)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m (_MOD (_MUL r a) m))))
`)

// POWMOD := Y (λrec.λa.λe.λm.IF (ISZERO e) (IF (ISZERO m) ONE (MOD ONE m)) (IF (ISEVEN e) (rec (MOD (MUL a a) m) (DIV2 e) m) (MOD (MUL a (rec (MOD (MUL a a) m) (DIV2 e) m)) m)))
var POWMOD = MakeLazyScript(`
	_Y (λrec.λa.λe.λm.
		_IF (_ISZERO e)
			(_IF (_ISZERO m) _ONE (_MOD _ONE m))
			(_IF (_ISEVEN e)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m)
				(_MOD (_MUL a (rec (_MOD (_MUL a a) m) (_DIV2 e) m)) m)))
`)

// Φ combinator for PRED
var (
	// Φ := λx.PAIR (SECOND x) (SUCC (SECOND x))
	PHI = MakeLazyScript(`λx._PAIR (_SECOND x) (_SUCC (_SECOND x))`)

	// PRED := λn.FIRST (n Φ (PAIR 0 0))
	PRED = MakeLazyScript(`λn._FIRST (n _PHI (_PAIR _0 _0))`)
)

// List operations
var (
	// NIL := λx.TRUE
	NIL = MakeLazyScript(`λx._TRUE`)

	// NULL := λp.p (λx.λy.FALSE)
	NULL = MakeLazyScript(`λp.p (λx.λy._FALSE)`)
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
var Y = MakeLazyScript(`λf.(λx.f (x x)) (λx.f (x x))`)

// FACTORIAL := Y (λf.λn.ISZERO n 1 (MULT n (f (PRED n))))
var FACTORIAL = MakeLazyScript(`
	_Y (λf.λn.
		(_ISZERO n) _1 (_MULT n (f (_PRED n))))
`)

// FAC is an alternative factorial implementation without Y combinator
// FAC = λn.λf.n(λf.λn.n(f(λf.λx.n f(f x))))(λx.f)(λx.x)
var FAC = MakeLazyScript(`λn.λf.n(λf.λn.n(f(λf.λx.n f(f x))))(λx.f)(λx.x)`)

// FIB is a Fibonacci implementation without Y combinator
// FIB = λn.λf.n(λc.λa.λb.c b(λx.a (b x)))(λx.λy.x)(λx.x)f
var FIB = MakeLazyScript(`λn.λf.n(λc.λa.λb.c b(λx.a (b x)))(λx.λy.x)(λx.x)f`)