package config

import (
	"testing"
)

func TestParseJSONC(t *testing.T) {
	input := `{
        // This is a comment
        "name": "test",
        "features": {
            "node": {}
        }
    }`

	result, err := ParseJSONC(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("expected name to be 'test', got %v", result["name"])
	}
}
