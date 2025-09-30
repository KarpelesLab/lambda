# lambdarun - Lambda Calculus Expression Evaluator

A command-line tool for evaluating lambda calculus expressions with Church encoding support.

## Installation

```bash
go build -o lambdarun
```

## Usage

```bash
lambdarun [options] <expression>
```

### Options

- `-steps int` - Maximum number of beta reduction steps (default: 10000)
- `-type string` - Output type: `auto`, `int`, `bool`, `lambda` (default: `auto`)

### Output Types

- **`auto`**: Automatically detects the result type (tries int first, then bool, then lambda)
- **`int`**: Forces interpretation as a Church numeral (integer)
- **`bool`**: Forces interpretation as a Church boolean
- **`lambda`**: Shows the raw lambda expression result

**Note:** Since `FALSE = ZERO` and `TRUE = ONE` in Church encoding (both are λf.λx. x and λf.λx. f x respectively),
the `auto` mode will show 0/1 for booleans. Use `-type bool` to explicitly show `true`/`false`.

## Examples

### Arithmetic

```bash
# Addition
$ lambdarun '_PLUS _2 _3'
5
Reduced in 6 steps

# Multiplication
$ lambdarun '_MULT _3 _4'
12
Reduced in 9 steps

# Exponentiation
$ lambdarun '_POW _2 _3'
8
Reduced in 16 steps
```

### Boolean Logic

```bash
# AND operation
$ lambdarun -type bool '_AND _TRUE _FALSE'
false
Reduced in 4 steps

# OR operation
$ lambdarun -type bool '_OR _TRUE _FALSE'
true
Reduced in 4 steps

# NOT operation
$ lambdarun -type bool '_NOT _TRUE'
false
Reduced in 3 steps
```

### Comparisons

```bash
# Less than or equal
$ lambdarun -type bool '_LEQ _2 _3'
true
Reduced in 60 steps

# Equal
$ lambdarun -type bool '_EQ _5 _5'
true
Reduced in 384 steps

# Less than
$ lambdarun -type bool '_LT _3 _2'
false
Reduced in 65 steps
```

### Custom Lambda Expressions

```bash
# Identity function
$ lambdarun '(\x. x) _5'
5
Reduced in 1 steps

# Show as lambda expression
$ lambdarun -type lambda '(\x. x) _5'
λf.λx.f (f (f (f (f x))))
Reduced in 1 steps
```

### Step Limit

```bash
# Set low step limit to see partial reduction
$ lambdarun -steps 5 '_MULT _3 _4'
λf.λx.f (f (f (f ((λf.λx.f (f (f (f x)))) f ((λf.λx.f (f (f (f x)))) f x)))))
Warning: Reached maximum step limit (5 steps)
Result may be partially reduced.
```

## Available Constants

### Church Numerals
- `_0`, `_1`, `_2`, `_3`, ... (any digit sequence)
- `_ZERO`, `_ONE`, `_TWO`

### Booleans
- `_TRUE`, `_FALSE`

### Arithmetic Operations
- `_SUCC` - Successor (n + 1)
- `_PLUS`, `_ADD` - Addition
- `_MULT` - Multiplication
- `_POW` - Exponentiation
- `_SUB` - Subtraction
- `_DEC` - Decrement (n - 1)
- `_MOD` - Modulo
- `_DIV2` - Integer division by 2

### Boolean Operations
- `_AND` - Logical AND
- `_OR` - Logical OR
- `_NOT` - Logical NOT
- `_IF` - If-then-else

### Comparison Operations
- `_LEQ` - Less than or equal (≤)
- `_LT` - Less than (<)
- `_EQ` - Equal (=)
- `_ISZERO` - Test if zero
- `_ISEVEN` - Test if even
- `_ISODD` - Test if odd

### Other Operations
- `_MAX` - Maximum of two numbers
- `_MIN` - Minimum of two numbers
- `_GCD` - Greatest common divisor

### Combinators
- `_I` - Identity
- `_K` - Constant
- `_S` - Substitution
- `_Y` - Y combinator (for recursion)

## Exit Codes

- `0` - Success
- `1` - Parse error, invalid arguments, or type mismatch

## Error Handling

### Parse Errors

```bash
$ lambdarun '(\\x. x'
Parse error: unbalanced parentheses: missing 1 closing parenthesis(es)
```

### Type Mismatch

```bash
$ lambdarun -type bool '_PLUS _2 _3'
Error: Result is not a valid Church boolean
λf.λx.f (f (f (f (f x))))
```

### Step Limit Reached

When the step limit is reached, a warning is printed to stderr:

```bash
Warning: Reached maximum step limit (5 steps)
Result may be partially reduced.
```

The partially reduced result is still printed to stdout.