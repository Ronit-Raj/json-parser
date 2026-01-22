package parser

import (
	"json_parser/scanner"
	"reflect"
	"testing"
)


// Helper function to reset scanner state before each test
func resetScanner() {
	scanner.Text = ""
	scanner.ResetPointer()
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		wantErr  bool
	}{
		{
			name:     "Simple string",
			input:    `"hello"`,
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "String with unicode",
			input:    `"hello world 🗿"`,
			expected: "hello world 🗿",
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    `""`,
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result string
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Integer",
			input:    `42`,
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "Negative integer",
			input:    `-123`,
			expected: -123,
			wantErr:  false,
		},
		{
			name:     "Float",
			input:    `3.14`,
			expected: 3.14,
			wantErr:  false,
		},
		{
			name:     "Negative float",
			input:    `-45.67`,
			expected: -45.67,
			wantErr:  false,
		},
		{
			name:     "Scientific notation",
			input:    `1e10`,
			expected: 1e10,
			wantErr:  false,
		},
		{
			name:     "Negative scientific notation",
			input:    `-2e3`,
			expected: -2000,
			wantErr:  false,
		},
		{
			name:     "Zero",
			input:    `0`,
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result float64
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeBoolAndNull(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		wantErr  bool
	}{
		{
			name:     "True",
			input:    `true`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "False",
			input:    `false`,
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result bool
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}

	t.Run("Null", func(t *testing.T) {
		resetScanner()
		var result any
		err := Decode(`null`, &result)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}
		if result != nil {
			t.Errorf("Decode() = %v, want nil", result)
		}
	})
	
}

func TestDecodeArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
		wantErr  bool
	}{		
		{
			name:     "Empty array",
			input:    `[]`,
			expected: []any{},
			wantErr:  false,
		},
		
		{
			name:     "Array of numbers",
			input:    `[1, 2, 3, 4, 5]`,
			expected: []any{float64(1), float64(2), float64(3), float64(4), float64(5)},
			wantErr:  false,
		},
		{
			name:     "Array of strings",
			input:    `["hello", "world"]`,
			expected: []any{"hello", "world"},
			wantErr:  false,
		},		
		{
			name:     "Mixed array",
			input:    `[1, "test", true, null, false]`,
			expected: []any{float64(1), "test", true, nil, false},
			wantErr:  false,
		},
		
		{
			name:     "Nested array",
			input:    `[1, [2, 3], 4]`,
			expected: []any{float64(1), []any{float64(2), float64(3)}, float64(4)},
			wantErr:  false,
		},
		{
			name:     "Array with whitespace",
			input:    `[ 1 , 2 , 3 ]`,
			expected: []any{float64(1), float64(2), float64(3)},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result []any
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
		wantErr  bool
	}{
		{
			name:     "Empty object",
			input:    `{}`,
			expected: map[string]any{},
			wantErr:  false,
		},
		{
			name:  "Simple object",
			input: `{"name": "John", "age": 30}`,
			expected: map[string]any{
				"name": "John",
				"age":  float64(30),
			},
			wantErr: false,
		},		
		{
			name:  "Object with different types",
			input: `{"string": "value", "number": 42, "bool": true, "null": null}`,
			expected: map[string]any{
				"string": "value",
				"number": float64(42),
				"bool":   true,
				"null":   nil,
			},
			wantErr: false,
		},
	
		{
			name:  "Nested object",
			input: `{"person": {"name": "Alice", "age": 25}}`,
			expected: map[string]any{
				"person": map[string]any{
					"name": "Alice",
					"age":  float64(25),
				},
			},
			wantErr: false,
		},
		{
			name:  "Object with array",
			input: `{"numbers": [1, 2, 3], "name": "test"}`,
			expected: map[string]any{
				"numbers": []any{float64(1), float64(2), float64(3)},
				"name":    "test",
			},
			wantErr: false,
		},
		{
			name:  "Object with whitespace",
			input: `{ "key1" : "value1" , "key2" : 123 }`,
			expected: map[string]any{
				"key1": "value1",
				"key2": float64(123),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result map[string]any
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeComplexStructures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
		wantErr  bool
	}{
		
			{
				name: "Complex nested structure",
				input: `{
					"class": 12,
					"sec": "A",
					"Name": "ronit🗿",
					"marks": {"phy": 90, "chem": 85, "maths": 90},
					"co-cirrcular": {},
					"address": null,
					"array": ["hello", "world"]
				}`,
				expected: map[string]any{
					"class": float64(12),
					"sec":   "A",
					"Name":  "ronit🗿",
					"marks": map[string]any{
						"phy":   float64(90),
						"chem":  float64(85),
						"maths": float64(90),
					},
					"co-cirrcular": map[string]any{},
					"address":      nil,
					"array":        []any{"hello", "world"},
				},
				wantErr: false,
			},
		
			{
				name:  "Array of objects",
				input: `{"users": [{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}]}`,
				expected: map[string]any{
					"users": []any{
						map[string]any{"name": "Alice", "age": float64(30)},
						map[string]any{"name": "Bob", "age": float64(25)},
					},
				},
				wantErr: false,
			},
			{
				name: "Deep nesting",
				input: `{
					"level1": {
						"level2": {
							"level3": {
								"value": "deep"
							}
						}
					}
				}`,
				expected: map[string]any{
					"level1": map[string]any{
						"level2": map[string]any{
							"level3": map[string]any{
								"value": "deep",
							},
						},
					},
				},
				wantErr: false,
			},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result map[string]any
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  any
		wantErr bool
	}{
		{
			name:    "Non-pointer target",
			input:   `{"key": "value"}`,
			target:  map[string]any{},
			wantErr: true,
		},
		{
			name:    "Type mismatch - string to number",
			input:   `"hello"`,
			target:  new(float64),
			wantErr: true,
		},
		{
			name:    "Type mismatch - object to array",
			input:   `{"key": "value"}`,
			target:  new([]any),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			err := Decode(tt.input, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeWithWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
		wantErr  bool
	}{
		{
			name: "Object with newlines and tabs",
			input: `{
				"key1": "value1",
				"key2": 123,
				"key3": [
					1,
					2,
					3
				]
			}`,
			expected: map[string]any{
				"key1": "value1",
				"key2": float64(123),
				"key3": []any{float64(1), float64(2), float64(3)},
			},
			wantErr: false,
		},
		{
			name:  "Extra whitespace everywhere",
			input: `  {  "key"  :  "value"  }  `,
			expected: map[string]any{
				"key": "value",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetScanner()
			var result map[string]any
			err := Decode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeEmptyStructures(t *testing.T) {
	t.Run("Empty object", func(t *testing.T) {
		resetScanner()
		var result map[string]any
		err := Decode(`{}`, &result)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %v", result)
		}
	})
	
	t.Run("Empty array", func(t *testing.T) {
		resetScanner()
		var result []any
		err := Decode(`[]`, &result)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}
		if len(result) != 0 {
			t.Errorf("Expected empty array, got %v", result)
		}
	})
	

	
	t.Run("Object with empty nested structures", func(t *testing.T) {
		resetScanner()
		var result map[string]any
		err := Decode(`{"obj": {}, "arr": []}`, &result)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}
		expected := map[string]any{
			"obj": map[string]any{},
			"arr": []any{},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Decode() = %v, want %v", result, expected)
		}
	})
	
}

// Benchmark tests
func BenchmarkDecodeSimpleObject(b *testing.B) {
	json := `{"name": "John", "age": 30, "city": "New York"}`
	for i := 0; i < b.N; i++ {
		resetScanner()
		var result map[string]any
		_ = Decode(json, &result)
	}
}


func BenchmarkDecodeComplexObject(b *testing.B) {
	json := `{
		"class": 12,
		"sec": "A",
		"Name": "ronit",
		"marks": {"phy": 90, "chem": 85, "maths": 90},
		"co-cirrcular": {},
		"address": null,
		"array": ["hello", "world"]
	}`
	for i := 0; i < b.N; i++ {
		resetScanner()
		var result map[string]any
		_ = Decode(json, &result)
	}
}


func BenchmarkDecodeArray(b *testing.B) {
	json := `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]`
	for i := 0; i < b.N; i++ {
		resetScanner()
		var result []any
		_ = Decode(json, &result)
	}
}

func BenchmarkDecodeNestedStructure(b *testing.B) {
	json := `{
		"level1": {
			"level2": {
				"level3": {
					"value": "deep",
					"array": [1, 2, 3]
				}
			}
		}
	}`
	for i := 0; i < b.N; i++ {
		resetScanner()
		var result map[string]any
		_ = Decode(json, &result)
	}
}
