package output

import "fmt"

// BuildLoggerAdapter adapts Output to builder.BuildLogger
type BuildLoggerAdapter struct {
	output Output
}

// NewBuildLoggerAdapter creates a new build logger adapter
func NewBuildLoggerAdapter(out Output) *BuildLoggerAdapter {
	return &BuildLoggerAdapter{output: out}
}

// Write writes a build log line
func (l *BuildLoggerAdapter) Write(line string) {
	if line != "" {
		l.output.Verbose(line)
	}
}

// buildLoggerWriter is a simple wrapper for writing build output
type buildLoggerWriter struct {
	output Output
}

// Write writes data to the output
func (w *buildLoggerWriter) Write(p []byte) (n int, err error) {
	w.output.Verbose(string(p))
	return len(p), nil
}

// NewBuildLoggerWriter creates a new build logger writer
func NewBuildLoggerWriter(out Output) *buildLoggerWriter {
	return &buildLoggerWriter{output: out}
}

// Printf prints a formatted message
func (w *buildLoggerWriter) Printf(format string, args ...interface{}) {
	w.output.Verbose(fmt.Sprintf(format, args...))
}
