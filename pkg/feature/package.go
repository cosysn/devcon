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
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func PackageFeature(dir string, output string) error {
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

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
	defer file.Close()

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
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Package the feature
	if err := PackageFeature(dir, tmpPath); err != nil {
		return fmt.Errorf("failed to package feature: %w", err)
	}

	// Load the tarball as an image
	img, err := tarball.ImageFromPath(tmpPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load tarball: %w", err)
	}

	// For now, push the tarball image directly
	// The feature metadata is embedded in the tarball files
	// (devcontainer-feature.json and install.sh)

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
