package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseJSONC parses JSON with Comments (JSONC) and returns a map
func ParseJSONC(input string) (map[string]interface{}, error) {
	// Remove comments and fix common JSONC issues
	cleanJSON, err := removeComments(input)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// removeComments removes single-line and multi-line comments from JSONC
func removeComments(input string) (string, error) {
	var result strings.Builder
	result.Grow(len(input))

	inString := false
	inLineComment := false
	inBlockComment := false
	prevChar := ' '

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		// Handle escape sequences in strings
		if inString {
			result.WriteRune(ch)
			if ch == '\\' && i+1 < len(input) {
				// Write the next character as part of the escape sequence
				i++
				result.WriteByte(input[i])
				prevChar = rune(input[i])
				continue
			}
			if ch == '"' {
				inString = false
			}
			prevChar = ch
			continue
		}

		// Start of string
		if ch == '"' {
			inString = true
			result.WriteRune(ch)
			prevChar = ch
			continue
		}

		// Line comment start: // found at position after previous char
		if !inBlockComment && i >= 1 && prevChar == '/' && ch == '/' {
			// Remove the previous '/' character by truncating output to exclude it
			current := result.String()
			result.Reset()
			result.Grow(len(current))
			result.WriteString(current[:len(current)-1])
			inLineComment = true
			prevChar = ch
			continue
		}

		// End of line comment
		if inLineComment && ch == '\n' {
			inLineComment = false
			result.WriteRune(ch)
			prevChar = ch
			continue
		}

		// Skip characters in line comment
		if inLineComment {
			prevChar = ch
			continue
		}

		// Block comment start: /* found
		if !inBlockComment && i >= 1 && prevChar == '/' && ch == '*' {
			// Remove the previous '/' character
			current := result.String()
			result.Reset()
			result.Grow(len(current))
			result.WriteString(current[:len(current)-1])
			inBlockComment = true
			prevChar = ch
			continue
		}

		// Block comment end: */
		if inBlockComment && prevChar == '*' && ch == '/' {
			inBlockComment = false
			prevChar = ch
			continue
		}

		// Skip characters in block comment
		if inBlockComment {
			prevChar = ch
			continue
		}

		result.WriteRune(ch)
		prevChar = ch
	}

	if inBlockComment || inLineComment {
		return "", fmt.Errorf("unclosed comment")
	}

	clean := result.String()

	// Handle trailing commas (e.g., {"a": 1,} -> {"a": 1})
	clean = fixTrailingCommas(clean)

	return clean, nil
}

// fixTrailingCommas removes trailing commas before closing braces/brackets
func fixTrailingCommas(input string) string {
	var result strings.Builder
	result.Grow(len(input))

	inString := false

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		// Handle escape sequences in strings
		if inString {
			result.WriteRune(ch)
			if ch == '\\' && i+1 < len(input) {
				i++
				result.WriteByte(input[i])
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		// Start of string
		if ch == '"' {
			inString = true
			result.WriteRune(ch)
			continue
		}

		// Trailing comma before } or ]
		if ch == ',' {
			// Look ahead to find next non-whitespace character
			nextIdx := i + 1
			skipComma := false
			for nextIdx < len(input) {
				nextCh := rune(input[nextIdx])
				if nextCh == ' ' || nextCh == '\t' || nextCh == '\n' || nextCh == '\r' {
					nextIdx++
					continue
				}
				if nextCh == '}' || nextCh == ']' {
					// Skip the comma
					skipComma = true
				}
				break
			}
			if skipComma {
				// Skip writing the comma
				continue
			}
		}

		result.WriteRune(ch)
	}

	return result.String()
}
