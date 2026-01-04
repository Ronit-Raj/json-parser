package scanner

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		wantErr  bool
	}{
		{
			name:  "Number Integer",
			input: "123",
			expected: []Token{
				{NumVal: 123, TypeOfToken: NUMBER},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Number Negative",
			input: "-45.6",
			expected: []Token{
				{NumVal: -45.6, TypeOfToken: NUMBER},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Number negative fraction",
			input: "-0.34",
			expected: []Token{
				{TypeOfToken: NUMBER, NumVal: -0.34},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Number negative fraction 02",
			input: "-2.34",
			expected: []Token{
				{TypeOfToken: NUMBER, NumVal: -2.34},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Number negative exponent",
			input: "-2e3",
			expected: []Token{
				{TypeOfToken: NUMBER, NumVal: -2000},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Number zero raise to exponent",
			input: "0e12",
			expected: []Token{
				{TypeOfToken: NUMBER, NumVal: 0},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "String Simple",
			input: `"hello"`,
			expected: []Token{
				{StringVal: "hello", TypeOfToken: STRING},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "String UTF8",
			input: `{"😅","😭"}`,
			expected: []Token{
				{TypeOfToken: BEGIN_OBJECT},
				{TypeOfToken: STRING, StringVal: `😅`},
				{TypeOfToken: VALUE_SEPARATOR},
				{TypeOfToken: STRING, StringVal: `😭`},
				{TypeOfToken: END_OBJECT},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Structural Tokens",
			input: `{}[]:,`,
			expected: []Token{
				{TypeOfToken: BEGIN_OBJECT},
				{TypeOfToken: END_OBJECT},
				{TypeOfToken: BEGIN_ARRAY},
				{TypeOfToken: END_ARRAY},
				{TypeOfToken: NAME_SEPARATOR},
				{TypeOfToken: VALUE_SEPARATOR},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Whitespace",
			input: "  \t\n 123 ",
			expected: []Token{
				{NumVal: 123, TypeOfToken: NUMBER},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Whitespaces02",
			input: "\t\t[\n\n]",
			expected: []Token{
				{TypeOfToken: BEGIN_ARRAY},
				{TypeOfToken: END_ARRAY},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "Complex JSON",
			input: `{"key": 123, "list": [1, 2]}`,
			expected: []Token{
				{TypeOfToken: BEGIN_OBJECT},
				{StringVal: "key", TypeOfToken: STRING},
				{TypeOfToken: NAME_SEPARATOR},
				{NumVal: 123, TypeOfToken: NUMBER},
				{TypeOfToken: VALUE_SEPARATOR},
				{StringVal: "list", TypeOfToken: STRING},
				{TypeOfToken: NAME_SEPARATOR},
				{TypeOfToken: BEGIN_ARRAY},
				{NumVal: 1, TypeOfToken: NUMBER},
				{TypeOfToken: VALUE_SEPARATOR},
				{NumVal: 2, TypeOfToken: NUMBER},
				{TypeOfToken: END_ARRAY},
				{TypeOfToken: END_OBJECT},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "CustomTest01",
			input: `-12.3e1[]`,
			expected: []Token{
				{TypeOfToken: NUMBER, NumVal: -123},
				{TypeOfToken: BEGIN_ARRAY},
				{TypeOfToken: END_ARRAY},
				{TypeOfToken: EOF},
			},
		},
		{
			name:    "CustomTest02",
			input:   "{@@@}",
			wantErr: true,
		},
		{
			name:    "CustomTest03",
			input:   "{😅😅}",
			wantErr: true,
		},
		{
			name:  "CustomTest04",
			input: `"{}"`,
			expected: []Token{
				{TypeOfToken: STRING, StringVal: `{}`},
				{TypeOfToken: EOF},
			},
		},
		{
			name:    "CustomTest05",
			input:   `"""`,
			wantErr: true,
		},
		{
			name:    "CustomTest06",
			input:   `12.{}`,
			wantErr: true,
		},
		{
			name:    "invalid numbers",
			input:   `12.`,
			wantErr: true,
		},
		{
			name:    "invalid numbers 02",
			input:   `001`,
			expected: []Token{
				{TypeOfToken: NUMBER,NumVal: 0},
				{TypeOfToken: NUMBER,NumVal: 0},
				{TypeOfToken: NUMBER,NumVal: 1},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "literal 01",
			input: `true`,
			expected: []Token{
				{TypeOfToken: LITERAL_TRUE},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "literal 02",
			input: `false`,
			expected: []Token{
				{TypeOfToken: LITERAL_FALSE},
				{TypeOfToken: EOF},
			},
		},
		{
			name:  "literal 03",
			input: `null`,
			expected: []Token{
				{TypeOfToken: LITERAL_NULL},
				{TypeOfToken: EOF},
			},
		},
		{
			name:    "invalid literal",
			input:   `nulls`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Text = tt.input
			pointer = 0

			var allTokens []Token
			var lastErr Error

			// Collect all tokens
			for {
				got, err := NextToken()
				allTokens = append(allTokens, got)
				lastErr = err

				if err.Code != 0 {
					break
				}
				if got.TypeOfToken == EOF {
					break
				}
			}

			// Check if error expectation matches
			if tt.wantErr {
				if lastErr.Code == 0 {
					t.Errorf("expected error but got none")
				}
				return
			}

			if lastErr.Code != 0 {
				t.Errorf("unexpected error: %v, msg: %s", lastErr, lastErr.Msg)
				return
			}

			// Check exact token count
			if len(allTokens) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d tokens", len(tt.expected), len(allTokens))
				t.Errorf("Expected tokens: %+v", tt.expected)
				t.Errorf("Got tokens: %+v", allTokens)
				return
			}

			// Verify each token
			for i, want := range tt.expected {
				got := allTokens[i]

				if got.TypeOfToken != want.TypeOfToken {
					t.Errorf("token %d: expected token type %v, got %v", i, want.TypeOfToken, got.TypeOfToken)
				}

				if want.TypeOfToken == NUMBER {
					if got.NumVal != want.NumVal {
						t.Errorf("token %d: expected number %v, got %v", i, want.NumVal, got.NumVal)
					}
				}

				if want.TypeOfToken == STRING {
					if got.StringVal != want.StringVal {
						t.Errorf("token %d: expected string %q, got %q", i, want.StringVal, got.StringVal)
					}
				}
			}
		})
	}
}
