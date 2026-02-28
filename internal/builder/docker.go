package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerBuilder struct {
	client *client.Client
}

func NewDockerBuilder() (*DockerBuilder, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
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
	contextPath := "."
	dockerfilePath := "Dockerfile"

	if spec.Dockerfile != "" {
		contextPath = filepath.Dir(spec.Dockerfile)
		dockerfilePath = filepath.Base(spec.Dockerfile)
	}

	// Read dockerfile content
	_, err := os.ReadFile(filepath.Join(contextPath, dockerfilePath))
	if err != nil {
		return "", err
	}

	// Build
	buildContext, err := os.Open(contextPath)
	if err != nil {
		return "", err
	}
	defer buildContext.Close()

	resp, err := b.client.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       []string{"devcon-build:latest"},
		Remove:     true,
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Note: In production, you'd parse the response to get the image ID
	// For now, return the tag
	return "devcon-build:latest", nil
}

func (b *DockerBuilder) Up(ctx context.Context, spec Spec) error {
	// Create and start container
	resp, err := b.client.ContainerCreate(ctx, &container.Config{
		Image: spec.Image,
		Env:   envToSlice(spec.Env),
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
