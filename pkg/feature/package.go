package feature

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func PackageFeature(dir string, output string) error {
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close outFile: %v\n", err)
		}
	}()

	gw := gzip.NewWriter(outFile)
	defer func() {
		if err := gw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close gzip writer: %v\n", err)
		}
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close tar writer: %v\n", err)
		}
	}()

	files := []string{
		"devcontainer-feature.json",
		"install.sh",
	}

	for _, file := range files {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		if err := addFileToTar(tw, path, file); err != nil {
			return err
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, path string, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", err)
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, name)
	if err != nil {
		return err
	}
	header.Name = name

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

// PublishFeature packages and publishes a Feature to OCI registry
func PublishFeature(ctx context.Context, dir string, ref string, opts ...remote.Option) error {
	// Create temp file for tarball
	tmpFile, err := os.CreateTemp("", "feature-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpPath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file: %v\n", err)
		}
	}()

	// Package the feature
	if err := PackageFeature(dir, tmpPath); err != nil {
		return fmt.Errorf("failed to package feature: %w", err)
	}

	// Load the tarball as a layer
	layer, err := tarball.LayerFromFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create layer from tarball: %w", err)
	}

	// Create an image from the layer using mutate.AppendLayers
	img, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return fmt.Errorf("failed to append layer to image: %w", err)
	}

	// Parse reference and push
	refParsed, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference: %w", err)
	}

	if err := remote.Write(refParsed, img, opts...); err != nil {
		return fmt.Errorf("failed to push to registry: %w", err)
	}

	return nil
}
