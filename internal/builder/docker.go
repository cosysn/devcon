package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
	"github.com/moby/go-archive/compression"
)

type DockerBuilder struct {
	client *client.Client
}

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
	_, err := os.ReadFile(dockerfileFullPath)
	if err != nil {
		return "", err
	}

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
	})
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close response body: %v\n", err)
		}
	}()

	// Read response body to ensure build completes
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read build response: %w", err)
	}

	// Note: In production, you'd parse the response to get the image ID
	// For now, return the tag
	return "devcon-build:latest", nil
}

func (b *DockerBuilder) Up(ctx context.Context, spec Spec) error {
	// Create and start container
	// Use /bin/sh as default command to keep container running
	cmd := []string{"/bin/sh", "-c", "while true; do sleep 3600; done"}
	resp, err := b.client.ContainerCreate(ctx, &container.Config{
		Image:        spec.Image,
		Env:          envToSlice(spec.Env),
		Cmd:          cmd,
		Tty:          true,
		AttachStdin:  true,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := b.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	fmt.Println("Container started:", resp.ID)
	return nil
}

func envToSlice(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}
