package main

import (
	"flag"
	"fmt"
	"os"

	lambda "github.com/KarpelesLab/lambda"
)

func main() {
	maxSteps := flag.Int("steps", 10000, "Maximum number of beta reduction steps")
	outputType := flag.String("type", "auto", "Output type: auto, int, bool, lambda")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <expression>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Evaluates a lambda calculus expression and prints the result.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s '_PLUS _2 _3'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type bool '_AND _TRUE _FALSE'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -steps 1000 '(\\x. x) _5'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type bool '_LEQ _2 _3'\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)

	// Parse the expression
	expr, err := lambda.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Reduce the expression
	result, steps := lambda.Reduce(expr, *maxSteps)

	// Check if we hit the step limit
	if steps >= *maxSteps {
		fmt.Fprintf(os.Stderr, "Warning: Reached maximum step limit (%d steps)\n", *maxSteps)
		fmt.Fprintf(os.Stderr, "Result may be partially reduced.\n\n")
	}

	// Handle output based on requested type
	switch *outputType {
	case "bool":
		// Force interpretation as boolean
		if b, ok := tryToBool(result); ok {
			fmt.Printf("%v\n", b)
		} else {
			fmt.Fprintf(os.Stderr, "Error: Result is not a valid Church boolean\n")
			fmt.Printf("%s\n", result)
			os.Exit(1)
		}

	case "int":
		// Force interpretation as integer
		if n, ok := tryToInt(result); ok {
			fmt.Printf("%d\n", n)
		} else {
			fmt.Fprintf(os.Stderr, "Error: Result is not a valid Church numeral\n")
			fmt.Printf("%s\n", result)
			os.Exit(1)
		}

	case "lambda":
		// Always show lambda expression
		fmt.Printf("%s\n", result)

	case "auto":
		// Try int first (since most operations produce numbers)
		if n, ok := tryToInt(result); ok {
			fmt.Printf("%d\n", n)
		} else if b, ok := tryToBool(result); ok {
			// Note: This won't be reached for 0/1 since they're valid ints
			fmt.Printf("%v\n", b)
		} else {
			// Show lambda expression
			fmt.Printf("%s\n", result)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid output type %q (must be: auto, int, bool, lambda)\n", *outputType)
		os.Exit(1)
	}

	if steps < *maxSteps {
		fmt.Fprintf(os.Stderr, "Reduced in %d steps\n", steps)
	}
}

// tryToInt attempts to interpret a Term as a Church numeral
// Returns the integer value and true if successful, or 0 and false otherwise
func tryToInt(obj lambda.Term) (int, bool) {
	// Church numerals have the form: λf.λx. f (f (f ... (f x)))
	// Try to extract the number of applications
	abs1, ok := obj.(lambda.Abstraction)
	if !ok {
		return 0, false
	}

	abs2, ok := abs1.Body.(lambda.Abstraction)
	if !ok {
		return 0, false
	}

	// Count the number of times the first parameter (f) is applied
	count := 0
	current := abs2.Body

	for {
		app, ok := current.(lambda.Application)
		if !ok {
			// Should end with the second parameter (x)
			if v, ok := current.(lambda.Var); ok && v.Name == abs2.Param {
				return count, true
			}
			return 0, false
		}

		// Check if function is the first parameter (f)
		if v, ok := app.Func.(lambda.Var); ok && v.Name == abs1.Param {
			count++
			current = app.Arg
		} else {
			return 0, false
		}
	}
}

// tryToBool attempts to interpret a Term as a Church boolean
// Returns the boolean value and true if successful, or false and false otherwise
func tryToBool(obj lambda.Term) (bool, bool) {
	// Church booleans have the form:
	// TRUE  = λx.λy.x
	// FALSE = λx.λy.y

	abs1, ok := obj.(lambda.Abstraction)
	if !ok {
		return false, false
	}

	abs2, ok := abs1.Body.(lambda.Abstraction)
	if !ok {
		return false, false
	}

	// The body should be a variable
	v, ok := abs2.Body.(lambda.Var)
	if !ok {
		return false, false
	}

	// Check which parameter it returns
	if v.Name == abs1.Param {
		return true, true
	} else if v.Name == abs2.Param {
		return false, true
	}

	return false, false
}