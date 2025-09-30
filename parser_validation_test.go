package lambda

import (
	"strings"
	"testing"
)

func TestParserParenthesesValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Balanced simple",
			input:   "(\\x. x)",
			wantErr: false,
		},
		{
			name:    "Balanced nested",
			input:   "((\\x. (y z)))",
			wantErr: false,
		},
		{
			name:    "Missing closing paren",
			input:   "(\\x. x",
			wantErr: true,
			errMsg:  "missing 1 closing parenthesis",
		},
		{
			name:    "Missing multiple closing parens",
			input:   "(((\\x. x)",
			wantErr: true,
			errMsg:  "missing 2 closing parenthesis",
		},
		{
			name:    "Extra closing paren",
			input:   "\\x. x)",
			wantErr: true,
			errMsg:  "unexpected ')'",
		},
		{
			name:    "Extra closing parens at end",
			input:   "(\\x. x))",
			wantErr: true,
			errMsg:  "unexpected ')'",
		},
		{
			name:    "Complex balanced expression",
			input:   "_IF (_LEQ _2 _3) (_OR (_EQ _2 _TWO) (_EQ _2 _3)) _FALSE",
			wantErr: false,
		},
		{
			name:    "Complex unbalanced expression",
			input:   "_IF (_LEQ _2 _3) (_OR (_EQ _2 _TWO) (_EQ _2 _3) _FALSE",
			wantErr: true,
			errMsg:  "missing 1 closing parenthesis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error containing %q, got nil", tt.input, tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Parse(%q) expected error containing %q, got %q", tt.input, tt.errMsg, err.Error())
				} else {
					t.Logf("✓ Correctly detected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				} else {
					t.Logf("✓ Parsed successfully")
				}
			}
		})
	}
}