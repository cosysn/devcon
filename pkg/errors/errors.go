package errors

import (
	"fmt"
)

// ErrorCode represents a category of errors
type ErrorCode string

const (
	// Config errors
	ErrCodeConfigNotFound     ErrorCode = "CONFIG_NOT_FOUND"
	ErrCodeConfigParse        ErrorCode = "CONFIG_PARSE_ERROR"
	ErrCodeConfigInvalid      ErrorCode = "CONFIG_INVALID"
	ErrCodeConfigExtend       ErrorCode = "CONFIG_EXTEND_ERROR"

	// Build errors
	ErrCodeBuildFailed    ErrorCode = "BUILD_FAILED"
	ErrCodeBuilderInit    ErrorCode = "BUILDER_INIT_FAILED"
	ErrCodeImageNotFound  ErrorCode = "IMAGE_NOT_FOUND"

	// Container errors
	ErrCodeContainerStart ErrorCode = "CONTAINER_START_FAILED"
	ErrCodeContainerCreate ErrorCode = "CONTAINER_CREATE_FAILED"

	// Docker errors
	ErrCodeDockerNotAvailable ErrorCode = "DOCKER_NOT_AVAILABLE"
	ErrCodeDockerConnection    ErrorCode = "DOCKER_CONNECTION_FAILED"

	// Feature errors
	ErrCodeFeatureResolve    ErrorCode = "FEATURE_RESOLVE_FAILED"
	ErrCodeFeatureDownload   ErrorCode = "FEATURE_DOWNLOAD_FAILED"
	ErrCodeFeaturePackage    ErrorCode = "FEATURE_PACKAGE_FAILED"

	// Generic errors
	ErrCodeUnknown ErrorCode = "UNKNOWN_ERROR"
)

// Suggestion provides actionable advice for fixing an error
type Suggestion struct {
	Text   string
	Action string
}

// EnhancedError contains error details with code and suggestions
type EnhancedError struct {
	Code       ErrorCode
	Message    string
	Suggestion *Suggestion
	Err       error
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *EnhancedError) Unwrap() error {
	return e.Err
}

// NewEnhancedError creates a new enhanced error
func NewEnhancedError(code ErrorCode, message string, err error, suggestion *Suggestion) *EnhancedError {
	return &EnhancedError{
		Code:       code,
		Message:    message,
		Suggestion: suggestion,
		Err:        err,
	}
}

// WithSuggestion adds a suggestion to an error
func WithSuggestion(err error, text, action string) *EnhancedError {
	if e, ok := err.(*EnhancedError); ok {
		e.Suggestion = &Suggestion{Text: text, Action: action}
		return e
	}
	return &EnhancedError{
		Code:       ErrCodeUnknown,
		Message:    err.Error(),
		Suggestion: &Suggestion{Text: text, Action: action},
		Err:        err,
	}
}

// Config error helpers

// NewConfigNotFoundError creates an error for missing config file
func NewConfigNotFoundError(path string) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeConfigNotFound,
		Message: fmt.Sprintf("devcontainer.json not found at %s", path),
		Suggestion: &Suggestion{
			Text:   "Create a .devcontainer/devcontainer.json file",
			Action: "See https://containers.dev/implementors/json_reference for the schema",
		},
	}
}

// NewConfigParseError creates an error for config parsing failures
func NewConfigParseError(path string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeConfigParse,
		Message: fmt.Sprintf("Failed to parse %s: %v", path, err),
		Suggestion: &Suggestion{
			Text:   "Check your devcontainer.json for syntax errors",
			Action: "Use a JSON validator or IDE with JSON schema support",
		},
		Err: err,
	}
}

// NewMissingImageOrDockerfileError creates an error for missing image/dockerfile
func NewMissingImageOrDockerfileError() *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeConfigInvalid,
		Message: "Either 'image' or 'dockerFile' must be specified in devcontainer.json",
		Suggestion: &Suggestion{
			Text:   "Add either an 'image' or 'dockerFile' property to your devcontainer.json",
			Action: "Example: {\"image\": \"mcr.microsoft.com/devcontainers/base:ubuntu\"}",
		},
	}
}

// Docker error helpers

// NewDockerNotAvailableError creates an error when Docker is not available
func NewDockerNotAvailableError() *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeDockerNotAvailable,
		Message: "Docker is not available or not running",
		Suggestion: &Suggestion{
			Text:   "Make sure Docker is installed and running",
			Action: "Run 'docker ps' to verify Docker is accessible",
		},
	}
}

// NewDockerConnectionError creates an error when Docker connection fails
func NewDockerConnectionError(err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeDockerConnection,
		Message: fmt.Sprintf("Failed to connect to Docker: %v", err),
		Suggestion: &Suggestion{
			Text:   "Check your Docker socket and permissions",
			Action: "You may need to add your user to the 'docker' group",
		},
		Err: err,
	}
}

// Feature error helpers

// NewFeatureResolveError creates an error when feature resolution fails
func NewFeatureResolveError(featureID string, err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeFeatureResolve,
		Message: fmt.Sprintf("Failed to resolve feature '%s': %v", featureID, err),
		Suggestion: &Suggestion{
			Text:   "Check the feature ID and version",
			Action: "Ensure the feature exists in the specified registry",
		},
		Err: err,
	}
}

// Build error helpers

// NewBuildFailedError creates an error when build fails
func NewBuildFailedError(err error) *EnhancedError {
	return &EnhancedError{
		Code:    ErrCodeBuildFailed,
		Message: fmt.Sprintf("Build failed: %v", err),
		Suggestion: &Suggestion{
			Text:   "Check the Dockerfile for errors",
			Action: "Run 'docker build' manually to see detailed error output",
		},
		Err: err,
	}
}
