package config

import (
	"encoding/json"

	"github.com/tailscale/hujson"
)

func ParseJSONC(input string) (map[string]interface{}, error) {
	ast, err := hujson.Parse([]byte(input))
	if err != nil {
		return nil, err
	}

	// Standardize converts HuJSON to standard JSON (strips comments, adds commas)
	// It mutates the AST in place
	ast.Standardize()

	// Pack the AST to get standard JSON bytes
	stdJSON := ast.Pack()

	var result map[string]interface{}
	if err := json.Unmarshal(stdJSON, &result); err != nil {
		return nil, err
	}

	return result, nil
}
