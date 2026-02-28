package config

import (
	"encoding/json"
	"testing"
)

func TestParseJSONC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantKey  string
		wantVal  interface{}
		wantErr  bool
	}{
		{
			name:    "basic object with comment",
			input:   `{"name": "test"}`,
			wantKey: "name",
			wantVal: "test",
		},
		{
			name:    "with line comment",
			input:   "{\n// comment\n\"a\": 1}",
			wantKey: "a",
			wantVal: float64(1),
		},
		{
			name:    "with block comment",
			input:   `{"a": 1 /* block */}`,
			wantKey: "a",
			wantVal: float64(1),
		},
		{
			name:    "trailing comma",
			input:   `{"a": 1,}`,
			wantKey: "a",
			wantVal: float64(1),
		},
		{
			name:    "nested object",
			input:   `{"nested": {"key": "value"}}`,
			wantKey: "nested",
			wantVal: map[string]interface{}{"key": "value"},
		},
		{
			name:    "array",
			input:   `{"arr": [1, 2, 3]}`,
			wantKey: "arr",
			wantVal: []interface{}{float64(1), float64(2), float64(3)},
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSONC(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := result[tt.wantKey]
			// Use JSON marshal for comparison since maps cannot be compared directly
			gotJSON, _ := json.Marshal(got)
			wantJSON, _ := json.Marshal(tt.wantVal)
			if string(gotJSON) != string(wantJSON) {
				t.Errorf("got %v, want %v", got, tt.wantVal)
			}
		})
	}
}
