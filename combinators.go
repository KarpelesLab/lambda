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
	// IF := λb.λx.λy.b x y
	IF = Abstraction{
		Param: "b",
		Body: Abstraction{
			Param: "x",
			Body: Abstraction{
				Param: "y",
				Body: Application{
					Func: Application{
						Func: Var{Name: "b"},
						Arg:  Var{Name: "x"},
					},
					Arg: Var{Name: "y"},
				},
			},
		},
	}

	// IFTHENELSE := λp.λa.λb.p a b (same as IF)
	IFTHENELSE = IF
)

// Arithmetic operations
var (
	// ZERO := λf.λx.x (Church numeral 0)
	ZERO = Abstraction{
		Param: "f",
		Body: Abstraction{
			Param: "x",
			Body:  Var{Name: "x"},
		},
	}

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

	// ONE := SUCC ZERO (Church numeral 1)
	ONE = Application{
		Func: SUCC,
		Arg:  ZERO,
	}

	// TWO := SUCC ONE (Church numeral 2)
	TWO = Application{
		Func: SUCC,
		Arg:  ONE,
	}

	// THREE := SUCC TWO (Church numeral 3)
	THREE = Application{
		Func: SUCC,
		Arg:  TWO,
	}

	// DEC := PRED (decrement, alias for predecessor)
	DEC = PRED

	// ADD := PLUS (addition, alias)
	ADD = PLUS

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
					Arg:  PRED,
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

	// POW := λb.λn.n b (exponentiation: b^n)
	POW = Abstraction{
		Param: "b",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: Var{Name: "n"},
				Arg:  Var{Name: "b"},
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

	// LT := λm.λn.NOT (LEQ n m)
	LT = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: NOT,
				Arg: Application{
					Func: Application{
						Func: LEQ,
						Arg:  Var{Name: "n"},
					},
					Arg: Var{Name: "m"},
				},
			},
		},
	}

	// EQ := λm.λn.AND (LEQ m n) (LEQ n m)
	EQ = Abstraction{
		Param: "m",
		Body: Abstraction{
			Param: "n",
			Body: Application{
				Func: Application{
					Func: AND,
					Arg: Application{
						Func: Application{
							Func: LEQ,
							Arg:  Var{Name: "m"},
						},
						Arg: Var{Name: "n"},
					},
				},
				Arg: Application{
					Func: Application{
						Func: LEQ,
						Arg:  Var{Name: "n"},
					},
					Arg: Var{Name: "m"},
				},
			},
		},
	}

	// MAX := λa.λb.IF (LEQ a b) b a
	MAX = Abstraction{
		Param: "a",
		Body: Abstraction{
			Param: "b",
			Body: Application{
				Func: Application{
					Func: Application{
						Func: IF,
						Arg: Application{
							Func: Application{
								Func: LEQ,
								Arg:  Var{Name: "a"},
							},
							Arg: Var{Name: "b"},
						},
					},
					Arg: Var{Name: "b"},
				},
				Arg: Var{Name: "a"},
			},
		},
	}

	// MIN := λa.λb.IF (LEQ a b) a b
	MIN = Abstraction{
		Param: "a",
		Body: Abstraction{
			Param: "b",
			Body: Application{
				Func: Application{
					Func: Application{
						Func: IF,
						Arg: Application{
							Func: Application{
								Func: LEQ,
								Arg:  Var{Name: "a"},
							},
							Arg: Var{Name: "b"},
						},
					},
					Arg: Var{Name: "a"},
				},
				Arg: Var{Name: "b"},
			},
		},
	}
)

// GCD := Y (λrec.λa.λb.IF (ISZERO b) a (rec b (MOD a b)))
var GCD = MakeLazyScript(`
	_Y (\rec.\a.\b.
		_IF (_ISZERO b)
			a
			(rec b (_MOD a b)))
`)

// MOD := Y (λrec.λm.λn.(ISZERO n) ZERO ((LT m n) m (rec (SUB m n) n)))
// Modulo operation with zero-divisor guard: m mod n = 0 if n = 0
var MOD = MakeLazyScript(`
	_Y (\rec.\m.\n.
		(_ISZERO n) _ZERO
		((_LT m n) m (rec (_SUB m n) n)))
`)

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

// Bit manipulation helpers
var (
	// STEP2 := λp.PAIR (IF (SECOND p) (SUCC (FIRST p)) (FIRST p)) (NOT (SECOND p))
	STEP2 = Abstraction{
		Param: "p",
		Body: Application{
			Func: Application{
				Func: PAIR,
				Arg: Application{
					Func: Application{
						Func: Application{
							Func: IF,
							Arg: Application{
								Func: SECOND,
								Arg:  Var{Name: "p"},
							},
						},
						Arg: Application{
							Func: SUCC,
							Arg: Application{
								Func: FIRST,
								Arg:  Var{Name: "p"},
							},
						},
					},
					Arg: Application{
						Func: FIRST,
						Arg:  Var{Name: "p"},
					},
				},
			},
			Arg: Application{
				Func: NOT,
				Arg: Application{
					Func: SECOND,
					Arg:  Var{Name: "p"},
				},
			},
		},
	}

	// INIT2 := PAIR ZERO FALSE
	INIT2 = Application{
		Func: Application{
			Func: PAIR,
			Arg:  ZERO,
		},
		Arg: FALSE,
	}
)

// Division and parity operations
var (
	// DIV2 := λn.FIRST (n STEP2 INIT2)
	DIV2 = Abstraction{
		Param: "n",
		Body: Application{
			Func: FIRST,
			Arg: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg:  STEP2,
				},
				Arg: INIT2,
			},
		},
	}

	// ISODD := λn.SECOND (n STEP2 INIT2)
	ISODD = Abstraction{
		Param: "n",
		Body: Application{
			Func: SECOND,
			Arg: Application{
				Func: Application{
					Func: Var{Name: "n"},
					Arg:  STEP2,
				},
				Arg: INIT2,
			},
		},
	}

	// ISEVEN := λn.NOT (ISODD n)
	ISEVEN = Abstraction{
		Param: "n",
		Body: Application{
			Func: NOT,
			Arg: Application{
				Func: ISODD,
				Arg:  Var{Name: "n"},
			},
		},
	}
)

// MUL := λm.λn.λf.m (n f)
// Note: MUL is already defined above as MULT, but we need it for POWMOD
var MUL = MULT

// POWMOD' := Y (λrec.λa.λe.λm.λr.IF (ISZERO e) (IF (ISZERO m) r (MOD r m)) (IF (ISEVEN e) (rec (MOD (MUL a a) m) (DIV2 e) m r) (rec (MOD (MUL a a) m) (DIV2 e) m (MOD (MUL r a) m))))
// Tail-recursive modular exponentiation with accumulator
var POWMOD_PRIME = MakeLazyScript(`
	_Y (\rec.\a.\e.\m.\r.
		_IF (_ISZERO e)
			(_IF (_ISZERO m) r (_MOD r m))
			(_IF (_ISEVEN e)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m r)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m (_MOD (_MUL r a) m))))
`)

// POWMOD := Y (λrec.λa.λe.λm.IF (ISZERO e) (IF (ISZERO m) ONE (MOD ONE m)) (IF (ISEVEN e) (rec (MOD (MUL a a) m) (DIV2 e) m) (MOD (MUL a (rec (MOD (MUL a a) m) (DIV2 e) m)) m)))
var POWMOD = MakeLazyScript(`
	_Y (\rec.\a.\e.\m.
		_IF (_ISZERO e)
			(_IF (_ISZERO m) _ONE (_MOD _ONE m))
			(_IF (_ISEVEN e)
				(rec (_MOD (_MUL a a) m) (_DIV2 e) m)
				(_MOD (_MUL a (rec (_MOD (_MUL a a) m) (_DIV2 e) m)) m)))
`)

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
var FACTORIAL = MakeLazyScript(`
	_Y (\f.\n.
		(_ISZERO n) _1 (_MULT n (f (_PRED n))))
`)

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