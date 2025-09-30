package lambda

// DECOMPOSE and related functions for Miller-Rabin
var (
	// TWODEC := Y (λrec.λs.λd.IF (ISEVEN d) (rec (SUCC s) (DIV2 d)) (PAIR s d))
	TWODEC = MakeLazyScript(`
		_Y (\rec.\s.\d.
			_IF (_ISEVEN d)
				(rec (_SUCC s) (_DIV2 d))
				(_PAIR s d))
	`)

	// DECOMPOSE := λn.TWODEC ZERO (DEC n)
	DECOMPOSE = MakeLazyScript(`
		\n. (_Y (\rec.\s.\d.
			_IF (_ISEVEN d)
				(rec (_SUCC s) (_DIV2 d))
				(_PAIR s d))) _ZERO (_DEC n)
	`)

	// Helper for LET expressions (λx.λf.f x)
	LET = MakeLazyScript(`\x.\f. f x`)

	// OR already exists, but we need it here
	OR_EXPR = OR

	// IS_LESS2 := λn.LEQ n (SUCC (SUCC ZERO))
	IS_LESS2 = MakeLazyScript(`\n. _LEQ n (_SUCC (_SUCC _ZERO))`)

	// IS_SMALL := λn.LEQ n (SUCC (SUCC (SUCC ZERO)))
	IS_SMALL = MakeLazyScript(`\n. _LEQ n (_SUCC (_SUCC (_SUCC _ZERO)))`)

	// MR_PASS - Miller-Rabin single base check
	// This is complex, so let's build it step by step
	MR_PASS = MakeLazyScript(`
		\n.\a.
			(\x.\f. f x) (_Y (\rec.\s.\d.
				_IF (_ISEVEN d)
					(rec (_SUCC s) (_DIV2 d))
					(_PAIR s d)) _ZERO (_DEC n)) (\sd.
			(\x.\f. f x) (_FIRST sd) (\s.
			(\x.\f. f x) (_SECOND sd) (\d.
			(\x.\f. f x) (_IF (_LEQ n (_SUCC (_SUCC (_SUCC _ZERO)))) _TWO
				(_ADD (_SUCC (_SUCC _ZERO)) (_MOD a (_SUB n (_SUCC (_SUCC _ZERO)))))) (\abase.
			(\x.\f. f x) (_POWMOD_PRIME abase d n _ONE) (\x0.
				_IF (_OR (_EQ x0 _ONE) (_EQ x0 (_DEC n)))
					_TRUE
					((\x.\f. f x)
						(_Y (\loop.\j.\x.
							_IF (_ISZERO j)
								_FALSE
								((\x.\f. f x) (_MOD (_MUL x x) n) (\x2.
									_IF (_EQ x2 (_DEC n)) _TRUE (loop (_DEC j) x2)))))
						(\run. run (_DEC s) x0)))))))
	`)

	// MR_SCAN := Y (λrec.λn.λa.λlimit.IF (LT limit a) TRUE (IF (NOT (EQ (GCD n a) ONE)) FALSE (IF (MR_PASS n a) (rec n (SUCC a) limit) FALSE)))
	MR_SCAN = MakeLazyScript(`
		_Y (\rec.\n.\a.\limit.
			_IF (_LT limit a)
				_TRUE
				(_IF (_NOT (_EQ (_GCD n a) _ONE))
					_FALSE
					(_IF ((\n.\a.
						(\x.\f. f x) (_Y (\rec.\s.\d.
							_IF (_ISEVEN d)
								(rec (_SUCC s) (_DIV2 d))
								(_PAIR s d)) _ZERO (_DEC n)) (\sd.
						(\x.\f. f x) (_FIRST sd) (\s.
						(\x.\f. f x) (_SECOND sd) (\d.
						(\x.\f. f x) (_IF (_LEQ n (_SUCC (_SUCC (_SUCC _ZERO)))) _TWO
							(_ADD (_SUCC (_SUCC _ZERO)) (_MOD a (_SUB n (_SUCC (_SUCC _ZERO)))))) (\abase.
						(\x.\f. f x) (_POWMOD_PRIME abase d n _ONE) (\x0.
							_IF (_OR (_EQ x0 _ONE) (_EQ x0 (_DEC n)))
								_TRUE
								((\x.\f. f x)
									(_Y (\loop.\j.\x.
										_IF (_ISZERO j)
											_FALSE
											((\x.\f. f x) (_MOD (_MUL x x) n) (\x2.
												_IF (_EQ x2 (_DEC n)) _TRUE (loop (_DEC j) x2)))))
									(\run. run (_DEC s) x0))))))) n a)
						(rec n (_SUCC a) limit)
						_FALSE))))
	`)

	// IS_PRIME := λn.IF (IS_SMALL n) (OR (EQ n TWO) (EQ n (SUCC TWO))) (IF (ISEVEN n) FALSE ...)
	// For small n (n ≤ 3): return TRUE if n=2 or n=3
	// For even n > 3: return FALSE
	// For odd n > 3: run Miller-Rabin with bases 2..min(12, n-2)
	IS_PRIME = MakeLazyScript(`
		\n.
			_IF (_LEQ n _3)
				(_OR (_EQ n _TWO) (_EQ n _3))
				(_IF (_ISEVEN n)
					_FALSE
					((\x.\f. f x) (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC (_SUCC _ZERO)))))))))))) (\B.
						(_Y (\rec.\nn.\a.\limit.
							_IF (_LT limit a)
								_TRUE
								(_IF (_NOT (_EQ (_GCD nn a) _ONE))
									_FALSE
									(_IF ((\nn.\a.
										(\x.\f. f x) (_Y (\rec.\s.\d.
											_IF (_ISEVEN d)
												(rec (_SUCC s) (_DIV2 d))
												(_PAIR s d)) _ZERO (_DEC nn)) (\sd.
										(\x.\f. f x) (_FIRST sd) (\s.
										(\x.\f. f x) (_SECOND sd) (\d.
										(\x.\f. f x) (_IF (_LEQ nn (_SUCC (_SUCC (_SUCC _ZERO)))) _TWO
											(_ADD (_SUCC (_SUCC _ZERO)) (_MOD a (_SUB nn (_SUCC (_SUCC _ZERO)))))) (\abase.
										(\x.\f. f x) (_POWMOD_PRIME abase d nn _ONE) (\x0.
											_IF (_OR (_EQ x0 _ONE) (_EQ x0 (_DEC nn)))
												_TRUE
												((\x.\f. f x)
													(_Y (\loop.\j.\x.
														_IF (_ISZERO j)
															_FALSE
															((\x.\f. f x) (_MOD (_MUL x x) nn) (\x2.
																_IF (_EQ x2 (_DEC nn)) _TRUE (loop (_DEC j) x2)))))
													(\run. run (_DEC s) x0))))))) nn a)
										(rec nn (_SUCC a) limit)
										_FALSE)))) n _TWO (_MIN B (_DEC (_DEC n)))))))
	`)
)