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
			name:  "String Simple",
			input: `"hello"`,
			expected: []Token{
				{StringVal: "hello", TypeOfToken: STRING},
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
			name: "CustomTest01",
			input: `-12.3e1[]`,
			expected: []Token{
				{TypeOfToken: NUMBER,NumVal: -123},
				{TypeOfToken: BEGIN_ARRAY},
				{TypeOfToken: END_ARRAY},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Text = tt.input
			pointer = 0
			
			for i, want := range tt.expected {
				got, err := NextToken()
				if tt.wantErr {
					if err.Code == 0 {
						t.Errorf("step %d: expected error but got none", i)
					}
					return 
				} else if err.Code != 0 {
					t.Errorf("step %d: unexpected error: %v, msg: %s", i, err, err.Msg)
					return
				}

				if got.TypeOfToken != want.TypeOfToken {
					t.Errorf("step %d: expected token type %v, got %v", i, want.TypeOfToken, got.TypeOfToken)
				}
				
				if want.TypeOfToken == NUMBER {
					// Compare with epsilon if necessary, but for simple cases exact check might suffice
					if got.NumVal != want.NumVal {
						t.Errorf("step %d: expected number %v, got %v", i, want.NumVal, got.NumVal)
					}
				}
				
				if want.TypeOfToken == STRING {
					if got.StringVal != want.StringVal {
						t.Errorf("step %d: expected string %q, got %q", i, want.StringVal, got.StringVal)
					}
				}
			}
		})
	}
}
