package lambda

import (
	"fmt"
)

// Example demonstrating how to parse lambda expressions
func ExampleParse() {
	// Parse an identity function
	identity, _ := Parse("λx.x")
	fmt.Println(identity)
	// Output: λx.x
}

// Example parsing and reducing an application
func ExampleParse_application() {
	// Parse and reduce (λx.x) y → y
	expr, _ := Parse("(λx.x) y")
	reduced, _ := expr.BetaReduce()
	fmt.Println(reduced)
	// Output: y
}

// Example parsing Church numerals
func ExampleParse_churchNumeral() {
	// Parse Church numeral 2
	two, _ := Parse("λf.λx.f (f x)")
	fmt.Println("Parsed:", two)
	fmt.Println("As integer:", ToInt(two))
	// Output:
	// Parsed: λf.λx.f (f x)
	// As integer: 2
}

// Example parsing and evaluating arithmetic
func ExampleParse_arithmetic() {
	// Parse PLUS 1 2
	expr, _ := Parse("(λm.λn.λf.λx.m f (n f x)) (λf.λx.f x) (λf.λx.f (f x))")

	// Reduce to normal form
	result, _ := Reduce(expr, 1000)

	fmt.Println("Result:", ToInt(result))
	// Output: Result: 3
}

// Example using backslash notation
func ExampleParse_backslash() {
	// Backslash can be used instead of λ
	expr, _ := Parse("\\x.\\y.x")
	fmt.Println(expr)
	// Output: λx.λy.x
}