package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
	"github.com/moby/go-archive/compression"
)

// BuildLogger is an interface for receiving build progress updates
type BuildLogger interface {
	Write(line string)
}

// DockerBuilder implements the Builder interface for Docker
type DockerBuilder struct {
	client *client.Client
	logger BuildLogger
}

// NewDockerBuilder creates a new Docker builder
func NewDockerBuilder() (*DockerBuilder, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerBuilder{client: cli}, nil
}

// SetLogger sets the build logger
func (b *DockerBuilder) SetLogger(logger BuildLogger) {
	b.logger = logger
}

func (b *DockerBuilder) log(format string, args ...interface{}) {
	if b.logger != nil {
		b.logger.Write(fmt.Sprintf(format, args...))
	}
}

// Build builds a Docker image
func (b *DockerBuilder) Build(ctx context.Context, spec Spec) (string, error) {
	// If no Dockerfile is specified but we have an image, just return the image name
	if spec.Dockerfile == "" && spec.Image != "" {
		return spec.Image, nil
	}

	// Determine build context
	contextPath := spec.ContextDir
	if contextPath == "" {
		contextPath = "."
	}
	dockerfilePath := "Dockerfile"

	if spec.Dockerfile != "" {
		// If Dockerfile is an absolute path or starts with ., use as-is relative to ContextDir
		if filepath.IsAbs(spec.Dockerfile) || spec.Dockerfile[0] == '.' {
			contextPath = filepath.Join(spec.ContextDir, filepath.Dir(spec.Dockerfile))
			dockerfilePath = filepath.Base(spec.Dockerfile)
		} else {
			// Simple filename, use ContextDir
			contextPath = spec.ContextDir
			dockerfilePath = spec.Dockerfile
		}
	}

	// Read dockerfile content
	dockerfileFullPath := filepath.Join(contextPath, dockerfilePath)
	dockerfileContent, err := os.ReadFile(dockerfileFullPath)
	if err != nil {
		return "", err
	}

	b.log("Step 1/1 : FROM %s", extractBaseImage(string(dockerfileContent)))

	// Create build context as tar archive
	buildContext, err := archive.Tar(contextPath, compression.None)
	if err != nil {
		return "", fmt.Errorf("failed to create build context: %w", err)
	}
	defer func() {
		if err := buildContext.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close build context: %v\n", err)
		}
	}()

	// Use the new dockerbuild.ImageBuildOptions instead of deprecated types.ImageBuildOptions
	//nolint:staticcheck // SA1019: types.ImageBuildOptions is deprecated
	resp, err := b.client.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       []string{"devcon-build:latest"},
		Remove:     true,
		SuppressOutput: false,
	})
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close response body: %v\n", err)
		}
	}()

	// Read response body and stream output
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Try to parse as JSON and extract useful information
		var buildLog map[string]interface{}
		if err := json.Unmarshal([]byte(line), &buildLog); err == nil {
			// Extract stream message
			if stream, ok := buildLog["stream"]; ok {
				b.log("%s", strings.TrimSpace(stream.(string)))
			} else if aux, ok := buildLog["aux"]; ok {
				// This typically contains the final image ID
				b.log("Build completed: %v", aux)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading build output: %w", err)
	}

	// Note: In production, you'd parse the response to get the image ID
	// For now, return the tag
	return "devcon-build:latest", nil
}

// extractBaseImage extracts the base image from a Dockerfile
func extractBaseImage(dockerfile string) string {
	lines := strings.Split(dockerfile, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "FROM ") {
			return strings.TrimPrefix(line, "FROM ")
		}
	}
	return "unknown"
}

// Up starts a container from the built image
func (b *DockerBuilder) Up(ctx context.Context, spec Spec) error {
	// Execute onCreateCommand before container starts (during creation)
	if spec.OnCreateCommand != "" {
		fmt.Println("Executing onCreateCommand:", spec.OnCreateCommand)
		if err := b.execInContainer(ctx, spec.Image, spec.OnCreateCommand); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: onCreateCommand failed: %v\n", err)
		}
	}

	// Create and start container
	// Use /bin/sh as default command to keep container running
	containerCmd := []string{"/bin/sh", "-c", "while true; do sleep 3600; done"}
	resp, err := b.client.ContainerCreate(ctx, &container.Config{
		Image:        spec.Image,
		Env:          envToSlice(spec.Env),
		Cmd:          containerCmd,
		Tty:          true,
		AttachStdin:  true,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := b.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	// Execute postStartCommand after container starts
	if spec.PostStartCommand != "" {
		fmt.Println("Executing postStartCommand:", spec.PostStartCommand)
		if err := b.execInContainer(ctx, resp.ID, spec.PostStartCommand); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: postStartCommand failed: %v\n", err)
		}
	}

	// Execute postCreateCommand after container starts
	if spec.PostCreateCommand != "" {
		fmt.Println("Executing postCreateCommand:", spec.PostCreateCommand)
		if err := b.execInContainer(ctx, resp.ID, spec.PostCreateCommand); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: postCreateCommand failed: %v\n", err)
		}
	}

	fmt.Println("Container started:", resp.ID)
	return nil
}

// execInContainer executes a command inside a container
func (b *DockerBuilder) execInContainer(ctx context.Context, containerID string, command string) error {
	// Use docker exec
	cmd := exec.CommandContext(ctx, "docker", "exec", containerID, "sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func envToSlice(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}
