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