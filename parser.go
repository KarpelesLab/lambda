package lambda

import (
	"fmt"
	"strings"
	"unicode"
)

// Parser for lambda calculus expressions
type Parser struct {
	input string
	pos   int
}

// Parse parses a lambda expression string and returns the corresponding Object
// Supported syntax:
//   - Variables: x, y, foo, bar123
//   - Abstraction: λx.body or \x.body
//   - Application: f x or (f x)
//   - Parentheses for grouping: (expr)
func Parse(input string) (Object, error) {
	p := &Parser{input: strings.TrimSpace(input), pos: 0}
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Check if we've consumed all input
	p.skipWhitespace()
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected characters after expression at position %d: %q", p.pos, p.input[p.pos:])
	}

	return result, nil
}

// parseExpr parses a complete expression
func (p *Parser) parseExpr() (Object, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	// Check for lambda abstraction
	if p.peekRune() == 'λ' || p.peek() == '\\' {
		return p.parseAbstraction()
	}

	// Parse application (left-associative)
	return p.parseApplication()
}

// parseAbstraction parses a lambda abstraction: λx.body or \x.body
func (p *Parser) parseAbstraction() (Object, error) {
	// Consume lambda symbol
	if p.peekRune() == 'λ' {
		p.pos += len("λ") // λ is multi-byte UTF-8
	} else if p.peek() == '\\' {
		p.pos++
	} else {
		return nil, fmt.Errorf("expected λ or \\ at position %d", p.pos)
	}

	p.skipWhitespace()

	// Parse parameter name
	param := p.parseIdentifier()
	if param == "" {
		return nil, fmt.Errorf("expected parameter name at position %d", p.pos)
	}

	p.skipWhitespace()

	// Consume dot
	if p.peek() != '.' {
		return nil, fmt.Errorf("expected '.' after parameter at position %d", p.pos)
	}
	p.pos++

	// Parse body
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return Abstraction{Param: param, Body: body}, nil
}

// parseApplication parses function application (left-associative)
// Examples: f x, f x y (= (f x) y), (f x) y
func (p *Parser) parseApplication() (Object, error) {
	// Parse the first term
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	// Keep parsing terms and building left-associative applications
	for {
		p.skipWhitespace()

		// Check if we can parse another term
		if p.pos >= len(p.input) {
			break
		}

		// Stop if we see a closing paren
		if p.peek() == ')' {
			break
		}

		// Try to parse another term
		right, err := p.parseTerm()
		if err != nil {
			// Not an error, just no more terms
			break
		}

		// Build application (left-associative)
		left = Application{Func: left, Arg: right}
	}

	return left, nil
}

// parseTerm parses a single term (variable or parenthesized expression)
func (p *Parser) parseTerm() (Object, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	// Check for parenthesized expression
	if p.peek() == '(' {
		p.pos++
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		p.skipWhitespace()
		if p.peek() != ')' {
			return nil, fmt.Errorf("expected ')' at position %d", p.pos)
		}
		p.pos++

		return expr, nil
	}

	// Check for lambda abstraction
	if p.peekRune() == 'λ' || p.peek() == '\\' {
		return p.parseAbstraction()
	}

	// Parse variable
	name := p.parseIdentifier()
	if name == "" {
		return nil, fmt.Errorf("expected variable or '(' at position %d", p.pos)
	}

	return Var{Name: name}, nil
}

// parseIdentifier parses a variable name
func (p *Parser) parseIdentifier() string {
	start := p.pos

	// First character must be a letter or underscore
	if p.pos < len(p.input) {
		r := rune(p.input[p.pos])
		if unicode.IsLetter(r) || r == '_' {
			p.pos++
		} else {
			return ""
		}
	} else {
		return ""
	}

	// Subsequent characters can be letters, digits, or underscores
	for p.pos < len(p.input) {
		r := rune(p.input[p.pos])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			p.pos++
		} else {
			break
		}
	}

	return p.input[start:p.pos]
}

// skipWhitespace skips whitespace characters
func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

// peek returns the current character without consuming it (as byte)
func (p *Parser) peek() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

// peekRune returns the current character as a proper UTF-8 rune
func (p *Parser) peekRune() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	// Check if it's the lambda character (λ = U+03BB)
	if p.pos+1 < len(p.input) && p.input[p.pos] == 0xCE && p.input[p.pos+1] == 0xBB {
		return 'λ'
	}
	return rune(p.input[p.pos])
}