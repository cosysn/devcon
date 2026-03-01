package progress

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Spinner represents a progress spinner
type Spinner struct {
	message string
	stopped bool
	mu      sync.Mutex
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
	}
}

var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var doneChars = []string{"✓", "✓", "✓"}

// Start starts the spinner
func (s *Spinner) Start() {
	if !IsTTY() {
		fmt.Println(s.message + "...")
		return
	}

	go func() {
		idx := 0
		for {
			s.mu.Lock()
			if s.stopped {
				s.mu.Unlock()
				break
			}
			s.mu.Unlock()

			fmt.Printf("\r%s %s", spinnerChars[idx%len(spinnerChars)], s.message)
			idx++
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Stop stops the spinner with a message
func (s *Spinner) Stop(message string) {
	s.mu.Lock()
	s.stopped = true
	s.mu.Unlock()

	if !IsTTY() {
		fmt.Println(message)
		return
	}

	// Clear the line and show done
	fmt.Printf("\r%s %s\n", doneChars[0], message)
}

// IsTTY returns true if the output is a TTY
func IsTTY() bool {
	return isTerminal(os.Stdout)
}

func isTerminal(f *os.File) bool {
	// Check if stdout is a terminal
	// This is a simplified check - in production, use proper terminal detection
	return false
}

// SimpleProgress shows a simple progress message
func SimpleProgress(msg string) {
	if IsTTY() {
		fmt.Printf("\r%s", strings.Repeat(" ", 60))
		fmt.Printf("\r%s", msg)
	} else {
		fmt.Print(msg + "... ")
	}
}

// SimpleComplete shows completion message
func SimpleComplete(msg string) {
	if !IsTTY() {
		fmt.Println(msg)
	} else {
		fmt.Printf("\r✓ %s\n", msg)
	}
}
