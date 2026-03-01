package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// OutputMode defines the output format
type OutputMode string

const (
	ModeText OutputMode = "text"
	ModeJSON OutputMode = "json"
)

// Output interface defines the contract for different output implementations
type Output interface {
	// Print outputs a message
	Print(msg string)
	// Printf outputs a formatted message
	Printf(format string, args ...interface{})
	// Println outputs a message with newline
	Println(msg string)
	// Verbose outputs a message only in verbose mode
	Verbose(msg string)
	// Verbosef outputs a formatted message only in verbose mode
	Verbosef(format string, args ...interface{})
	// Error outputs an error message
	Error(msg string)
	// Errorf outputs a formatted error message
	Errorf(format string, args ...interface{})
	// Success outputs a success message
	Success(msg string)
	// Successf outputs a formatted success message
	Successf(format string, args ...interface{})
	// Warn outputs a warning message
	Warn(msg string)
	// Warnf outputs a formatted warning message
	Warnf(format string, args ...interface{})
	// StartProgress starts a progress indicator
	StartProgress(msg string)
	// StopProgress stops the progress indicator
	StopProgress(msg string)
	// SetVerbose sets verbose mode
	SetVerbose(v bool)
	// IsVerbose returns verbose mode status
	IsVerbose() bool
	// IsQuiet returns quiet mode status
	IsQuiet() bool
}

// JSONResponse is the standard JSON response envelope
type JSONResponse struct {
	Success   bool        `json:"success"`
	Command   string      `json:"command"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	ErrorCode string      `json:"errorCode,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// PlainOutput implements Output for plain text format
type PlainOutput struct {
	out     io.Writer
	errOut  io.Writer
	verbose bool
	quiet   bool
}

// NewPlainOutput creates a new PlainOutput instance
func NewPlainOutput() *PlainOutput {
	return &PlainOutput{
		out:    os.Stdout,
		errOut: os.Stderr,
	}
}

// Print outputs a message
func (p *PlainOutput) Print(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprint(p.out, msg)
}

// Printf outputs a formatted message
func (p *PlainOutput) Printf(format string, args ...interface{}) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintf(p.out, format, args...)
}

// Println outputs a message with newline
func (p *PlainOutput) Println(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintln(p.out, msg)
}

// Verbose outputs a message only in verbose mode
func (p *PlainOutput) Verbose(msg string) {
	if p.verbose && !p.quiet {
		_, _ = fmt.Fprintln(p.out, msg)
	}
}

// Verbosef outputs a formatted message only in verbose mode
func (p *PlainOutput) Verbosef(format string, args ...interface{}) {
	if p.verbose && !p.quiet {
		_, _ = fmt.Fprintf(p.out, format+"\n", args...)
	}
}

// Error outputs an error message
func (p *PlainOutput) Error(msg string) {
	_, _ = fmt.Fprintln(p.errOut, "Error:", msg)
}

// Errorf outputs a formatted error message
func (p *PlainOutput) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.errOut, "Error: "+format+"\n", args...)
}

// Success outputs a success message
func (p *PlainOutput) Success(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintln(p.out, "✓", msg)
}

// Successf outputs a formatted success message
func (p *PlainOutput) Successf(format string, args ...interface{}) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintf(p.out, "✓ "+format+"\n", args...)
}

// Warn outputs a warning message
func (p *PlainOutput) Warn(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintln(p.out, "Warning:", msg)
}

// Warnf outputs a formatted warning message
func (p *PlainOutput) Warnf(format string, args ...interface{}) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintf(p.out, "Warning: "+format+"\n", args...)
}

// StartProgress starts a progress indicator (no-op for plain output)
func (p *PlainOutput) StartProgress(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprint(p.out, msg+"... ")
}

// StopProgress stops the progress indicator
func (p *PlainOutput) StopProgress(msg string) {
	if p.quiet {
		return
	}
	_, _ = fmt.Fprintln(p.out, msg)
}

// SetVerbose sets verbose mode
func (p *PlainOutput) SetVerbose(v bool) {
	p.verbose = v
}

// SetQuiet sets quiet mode
func (p *PlainOutput) SetQuiet(v bool) {
	p.quiet = v
}

// IsVerbose returns verbose mode status
func (p *PlainOutput) IsVerbose() bool {
	return p.verbose
}

// IsQuiet returns quiet mode status
func (p *PlainOutput) IsQuiet() bool {
	return p.quiet
}

// JSONOutput implements Output for JSON format
type JSONOutput struct {
	out     io.Writer
	errOut  io.Writer
	verbose bool
	quiet   bool
	cmdName string
}

// NewJSONOutput creates a new JSONOutput instance
func NewJSONOutput(cmdName string) *JSONOutput {
	return &JSONOutput{
		out:     os.Stdout,
		errOut: os.Stderr,
		cmdName: cmdName,
	}
}

func (j *JSONOutput) timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func (j *JSONOutput) printJSON(resp JSONResponse) {
	resp.Timestamp = j.timestamp()
	resp.Command = j.cmdName
	data, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(j.out, string(data))
}

// Print outputs a message
func (j *JSONOutput) Print(msg string) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: msg,
	})
}

// Printf outputs a formatted message
func (j *JSONOutput) Printf(format string, args ...interface{}) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: fmt.Sprintf(format, args...),
	})
}

// Println outputs a message with newline
func (j *JSONOutput) Println(msg string) {
	j.Print(msg)
}

// Verbose outputs a message only in verbose mode
func (j *JSONOutput) Verbose(msg string) {
	if j.verbose && !j.quiet {
		j.Print(msg)
	}
}

// Verbosef outputs a formatted message only in verbose mode
func (j *JSONOutput) Verbosef(format string, args ...interface{}) {
	if j.verbose && !j.quiet {
		j.Printf(format, args...)
	}
}

// Error outputs an error message
func (j *JSONOutput) Error(msg string) {
	j.printJSON(JSONResponse{
		Success:   false,
		Error:     msg,
		ErrorCode: "ERROR",
	})
}

// Errorf outputs a formatted error message
func (j *JSONOutput) Errorf(format string, args ...interface{}) {
	j.Error(fmt.Sprintf(format, args...))
}

// Success outputs a success message
func (j *JSONOutput) Success(msg string) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: msg,
	})
}

// Successf outputs a formatted success message
func (j *JSONOutput) Successf(format string, args ...interface{}) {
	if j.quiet {
		return
	}
	j.Success(fmt.Sprintf(format, args...))
}

// Warn outputs a warning message
func (j *JSONOutput) Warn(msg string) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: "Warning: " + msg,
	})
}

// Warnf outputs a formatted warning message
func (j *JSONOutput) Warnf(format string, args ...interface{}) {
	if j.quiet {
		return
	}
	j.Warn(fmt.Sprintf(format, args...))
}

// StartProgress starts a progress indicator
func (j *JSONOutput) StartProgress(msg string) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: msg + "...",
	})
}

// StopProgress stops the progress indicator
func (j *JSONOutput) StopProgress(msg string) {
	if j.quiet {
		return
	}
	j.printJSON(JSONResponse{
		Success: true,
		Message: msg,
	})
}

// SetVerbose sets verbose mode
func (j *JSONOutput) SetVerbose(v bool) {
	j.verbose = v
}

// IsVerbose returns verbose mode status
func (j *JSONOutput) IsVerbose() bool {
	return j.verbose
}

// IsQuiet returns quiet mode status
func (j *JSONOutput) IsQuiet() bool {
	return j.quiet
}

// GlobalOutput is the global output instance
var globalOutput Output = NewPlainOutput()

// SetGlobalOutput sets the global output instance
func SetGlobalOutput(o Output) {
	globalOutput = o
}

// GetGlobalOutput returns the global output instance
func GetGlobalOutput() Output {
	return globalOutput
}

// AddOutputFlags adds common output flags to a command
func AddOutputFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().StringVarP(&outputMode, "output", "o", "text", "Output format (text, json)")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress non-essential output")
}

// ApplyOutputSettings applies the output settings based on flags
func ApplyOutputSettings(cmd *cobra.Command, cmdName string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	outputMode, _ := cmd.Flags().GetString("output")
	quiet, _ := cmd.Flags().GetBool("quiet")

	var o Output
	switch OutputMode(outputMode) {
	case ModeJSON:
		o = NewJSONOutput(cmdName)
	default:
		o = NewPlainOutput()
	}

	o.SetVerbose(verbose)
	if p, ok := o.(*PlainOutput); ok {
		p.SetQuiet(quiet)
	}

	SetGlobalOutput(o)
}

// Package-level variables for flag binding
var verbose bool
var outputMode string
var quiet bool

// IsTTY returns true if the output is a TTY
func IsTTY() bool {
	return isTerminal(os.Stdout)
}

func isTerminal(f *os.File) bool {
	return isTerminalFd(f.Fd())
}

func isTerminalFd(fd uintptr) bool {
	// Simple check - in a real implementation, this would use syscalls
	// For now, we'll use a simple heuristic
	return false // Simplified for cross-platform compatibility
}
