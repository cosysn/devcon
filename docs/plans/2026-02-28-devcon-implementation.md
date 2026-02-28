# Devcon Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现 Go 版 devcontainer CLI 工具 MVP，支持 WSL2/Linux 平台、Docker 后端、Feature 生命周期和 Devcontainer 构建。

**Architecture:** 四层架构 - CLI 层 + 配置解析层 + Feature 调度层 + 构建抽象层。先实现 Docker Builder，后续可扩展其他后端。

**Tech Stack:** Go, Cobra, go-containerregistry, docker/client, hujson

---

## 阶段一：项目初始化

### Task 1: 初始化 Go Module 和项目结构

**Files:**
- Create: `go.mod`
- Create: `cmd/devcon/main.go`
- Create: `cmd/devcon/root.go`
- Create: `internal/builder/interface.go`
- Create: `pkg/config/devcontainer.go`
- Create: `pkg/config/feature.go`

**Step 1: 初始化 Go module**

Run: `go mod init github.com/devcon/cli`
Expected: go.mod created

**Step 2: 创建项目目录结构**

```bash
mkdir -p cmd/devcon
mkdir -p internal/builder
mkdir -p internal/config
mkdir -p internal/feature
mkdir -p internal/registry
mkdir -p pkg/config
```

**Step 3: 创建空命令文件**

```go
// cmd/devcon/main.go
package main

import "github.com/devcon/cli/cmd/devcon"

func main() {
    cmd.Execute()
}
```

```go
// cmd/devcon/root.go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "devcon",
    Short: "devcontainer CLI tool",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(featuresCmd)
}
```

**Step 4: 添加 Cobra 依赖**

Run: `go get github.com/spf13/cobra`
Expected: dependency added

**Step 5: 创建 Builder Interface**

```go
// internal/builder/interface.go
package builder

import (
    "context"
)

type Spec struct {
    Image         string
    Dockerfile    string
    Features      map[string]interface{}
    Env           map[string]string
    Mounts        []string
    Ports         []int
    RemoteUser    string
}

type Builder interface {
    Build(ctx context.Context, spec Spec) (string, error)
    Up(ctx context.Context, spec Spec) error
}
```

**Step 6: Commit**

```bash
git add .
git commit -m "chore: initialize project structure"
```

---

### Task 2: 添加 features 子命令

**Files:**
- Modify: `cmd/devcon/root.go`

**Step 1: 添加 features 子命令**

```go
// cmd/devcon/root.go - 添加 featuresCmd
var featuresCmd = &cobra.Command{
    Use:   "features",
    Short: "Feature lifecycle management",
}

func init() {
    featuresCmd.AddCommand(featuresPackageCmd)
    featuresCmd.AddCommand(featuresPublishCmd)
    rootCmd.AddCommand(featuresCmd)
}
```

**Step 2: 添加 package 子命令**

```go
// cmd/devcon/features_package.go
package cmd

import (
    "github.com/spf13/cobra"
)

var featuresPackageCmd = &cobra.Command{
    Use:   "package <dir>",
    Short: "Package a Feature as OCI artifact",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}
```

**Step 3: 添加 publish 子命令**

```go
// cmd/devcon/features_publish.go
package cmd

import (
    "github.com/spf13/cobra"
)

var featuresPublishCmd = &cobra.Command{
    Use:   "publish <dir>",
    Short: "Publish a Feature to OCI registry",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}

func init() {
    featuresPublishCmd.Flags().String("reg", "", "Registry URL")
}
```

**Step 4: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: binary builds successfully

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add features package/publish commands"
```

---

### Task 3: 添加 build/config/inspect/up 子命令

**Files:**
- Modify: `cmd/devcon/root.go`

**Step 1: 添加 build 命令**

```go
// cmd/devcon/build.go
package cmd

import (
    "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
    Use:   "build <dir>",
    Short: "Build devcontainer image",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}

func init() {
    buildCmd.Flags().String("provider", "docker", "Build provider (docker)")
    rootCmd.AddCommand(buildCmd)
}
```

**Step 2: 添加 config 命令**

```go
// cmd/devcon/config.go
package cmd

import (
    "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
    Use:   "config <dir>",
    Short: "Validate and show devcontainer config",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}

func init() {
    rootCmd.AddCommand(configCmd)
}
```

**Step 3: 添加 inspect 命令**

```go
// cmd/devcon/inspect.go
package cmd

import (
    "github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
    Use:   "inspect <dir>",
    Short: "Inspect parsed config and feature dependencies",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}

func init() {
    rootCmd.AddCommand(inspectCmd)
}
```

**Step 4: 添加 up 命令**

```go
// cmd/devcon/up.go
package cmd

import (
    "github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
    Use:   "up <dir>",
    Short: "Start devcontainer for local testing",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil // TODO: implement
    },
}

func init() {
    rootCmd.AddCommand(upCmd)
}
```

**Step 5: 编译测试**

Run: `go build -o devcon ./cmd/devcon && ./devcon --help`
Expected: shows all commands

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add build/config/inspect/up commands"
```

---

## 阶段二：配置解析层

### Task 4: 实现 JSONC 解析

**Files:**
- Create: `pkg/config/jsonc.go`
- Create: `pkg/config/jsonc_test.go`

**Step 1: 创建测试**

```go
// pkg/config/jsonc_test.go
package config

import (
    "testing"
)

func TestParseJSONC(t *testing.T) {
    input := `{
        // This is a comment
        "name": "test",
        "features": {
            "node": {}
        }
    }`

    result, err := ParseJSONC(input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result["name"] != "test" {
        t.Errorf("expected name to be 'test', got %v", result["name"])
    }
}
```

**Step 2: 运行测试**

Run: `go test ./pkg/config/... -v`
Expected: FAIL - function not defined

**Step 3: 实现解析**

```go
// pkg/config/jsonc.go
package config

import (
    "github.com/tailscale/hujson"
)

func ParseJSONC(input string) (map[string]interface{}, error) {
    ast, err := hujson.Parse(input)
    if err != nil {
        return nil, err
    }
    ast.Normalize()

    var result map[string]interface{}
    if err := ast.Unmarshal(&result); err != nil {
        return nil, err
    }

    return result, nil
}
```

**Step 4: 运行测试**

Run: `go test ./pkg/config/... -v`
Expected: PASS

**Step 5: 添加 hujson 依赖**

Run: `go get github.com/tailscale/hujson`
Expected: dependency added

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add JSONC parsing support"
```

---

### Task 5: 实现 devcontainer.json 解析

**Files:**
- Modify: `pkg/config/devcontainer.go`
- Create: `pkg/config/devcontainer_test.go`

**Step 1: 创建类型定义**

```go
// pkg/config/devcontainer.go
package config

import (
    "os"
    "path/filepath"
)

type DevcontainerConfig struct {
    Image         string                 `json:"image,omitempty"`
    Dockerfile    string                 `json:"dockerFile,omitempty"`
    Build         *BuildConfig           `json:"build,omitempty"`
    Features      map[string]interface{} `json:"features,omitempty"`
    Env           map[string]string      `json:"containerEnv,omitempty"`
    Mounts        []string               `json:"mounts,omitempty"`
    Ports         []int                  `json:"forwardPorts,omitempty"`
    RemoteUser    string                 `json:"remoteUser,omitempty"`
    Extends       string                 `json:"extends,omitempty"`
}

type BuildConfig struct {
    Dockerfile string `json:"dockerfile,omitempty"`
    Context    string `json:"context,omitempty"`
}

func ParseDevcontainer(dir string) (*DevcontainerConfig, error) {
    path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    parsed, err := ParseJSONC(string(data))
    if err != nil {
        return nil, err
    }

    // Convert to struct
    // ... implementation
}
```

**Step 2: 创建测试**

```go
// pkg/config/devcontainer_test.go
package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParseDevcontainer(t *testing.T) {
    dir := t.TempDir()
    devcontainerDir := filepath.Join(dir, ".devcontainer")
    os.MkdirAll(devcontainerDir, 0755)

    config := `{
        "image": "mcr.microsoft.com/devcontainers/base:ubuntu",
        "features": {
            "node": {}
        }
    }`

    os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

    result, err := ParseDevcontainer(dir)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.Image != "mcr.microsoft.com/devcontainers/base:ubuntu" {
        t.Errorf("expected image, got %v", result.Image)
    }
}
```

**Step 3: 运行测试**

Run: `go test ./pkg/config/... -v`
Expected: FAIL - function incomplete

**Step 4: 完成实现**

```go
// pkg/config/devcontainer.go - 完善 ParseDevcontainer
func ParseDevcontainer(dir string) (*DevcontainerConfig, error) {
    path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    parsed, err := ParseJSONC(string(data))
    if err != nil {
        return nil, err
    }

    // Simple JSON unmarshal into struct
    jsonData, _ := json.Marshal(parsed)
    var config DevcontainerConfig
    if err := json.Unmarshal(jsonData, &config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

**Step 5: 添加 json 依赖**

Run: `go get encoding/json` (already in stdlib)`
Expected: no change

**Step 6: 运行测试**

Run: `go test ./pkg/config/... -v`
Expected: PASS

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add devcontainer.json parsing"
```

---

### Task 6: 实现 extends 继承处理

**Files:**
- Modify: `pkg/config/devcontainer.go`

**Step 1: 添加测试**

```go
// pkg/config/devcontainer_test.go - 添加 extends 测试
func TestParseDevcontainerWithExtends(t *testing.T) {
    dir := t.TempDir()
    devcontainerDir := filepath.Join(dir, ".devcontainer")
    os.MkdirAll(devcontainerDir, 0755)

    baseConfig := `{
        "image": "mcr.microsoft.com/devcontainers/base:ubuntu",
        "features": {
            "git": {}
        }
    }`
    os.WriteFile(filepath.Join(devcontainerDir, "base.json"), []byte(baseConfig), 0644)

    config := `{
        "extends": "./base.json",
        "features": {
            "node": {}
        }
    }`
    os.WriteFile(filepath.Join(devcontainerDir, "devcontainer.json"), []byte(config), 0644)

    result, err := ParseDevcontainer(dir)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Should have both git and node features
    if result.Features == nil {
        t.Fatal("features should not be nil")
    }
}
```

**Step 2: 运行测试**

Run: `go test ./pkg/config/... -v -run TestParseDevcontainerWithExtends`
Expected: FAIL - extends not implemented

**Step 3: 实现 extends**

```go
// pkg/config/devcontainer.go - 添加 ResolveExtends
func ResolveExtends(dir string, config *DevcontainerConfig) (*DevcontainerConfig, error) {
    if config.Extends == "" {
        return config, nil
    }

    // Resolve extends path
    extendsPath := filepath.Join(dir, ".devcontainer", config.Extends)
    data, err := os.ReadFile(extendsPath)
    if err != nil {
        return nil, err
    }

    parsed, err := ParseJSONC(string(data))
    if err != nil {
        return nil, err
    }

    var base DevcontainerConfig
    jsonData, _ := json.Marshal(parsed)
    json.Unmarshal(jsonData, &base)

    // Merge: base -> config (config takes precedence)
    if config.Image != "" {
        base.Image = config.Image
    }

    // Merge features
    if config.Features != nil {
        if base.Features == nil {
            base.Features = make(map[string]interface{})
        }
        for k, v := range config.Features {
            base.Features[k] = v
        }
    }

    // Merge env
    if config.Env != nil {
        if base.Env == nil {
            base.Env = make(map[string]string)
        }
        for k, v := range config.Env {
            base.Env[k] = v
        }
    }

    return &base, nil
}
```

**Step 4: 运行测试**

Run: `go test ./pkg/config/... -v -run TestParseDevcontainerWithExtends`
Expected: PASS

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add extends inheritance support"
```

---

## 阶段三：Feature 调度层

### Task 7: 实现 Feature 定义解析

**Files:**
- Create: `pkg/config/feature.go`
- Create: `pkg/config/feature_test.go`

**Step 1: 创建 Feature 类型**

```go
// pkg/config/feature.go
package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Feature struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name,omitempty"`
    Version     string                 `json:"version,omitempty"`
    Options     map[string]interface{} `json:"options,omitempty"`
    DependsOn   []string               `json:"dependsOn,omitempty"`
    ContainerEnv map[string]string    `json:"containerEnv,omitempty"`
}

type FeatureDefinition struct {
    ID        string                 `json:"id"`
    Name      string                 `json:"name"`
    Version   string                 `json:"version"`
    DependsOn []string               `json:"dependsOn"`
    Options   map[string]OptionSpec  `json:"options"`
}

type OptionSpec struct {
    Type        string      `json:"type"`
    Default     interface{} `json:"default"`
    Description string      `json:"description,omitempty"`
}

func ParseFeatureDefinition(dir string) (*FeatureDefinition, error) {
    path := filepath.Join(dir, "devcontainer-feature.json")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var def FeatureDefinition
    if err := json.Unmarshal(data, &def); err != nil {
        return nil, err
    }

    return &def, nil
}
```

**Step 2: 创建测试**

```go
// pkg/config/feature_test.go
package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParseFeatureDefinition(t *testing.T) {
    dir := t.TempDir()

    config := `{
        "id": "node",
        "name": "Node.js",
        "version": "1.0.0",
        "dependsOn": ["git"],
        "options": {
            "version": {
                "type": "string",
                "default": "20"
            }
        }
    }`

    os.WriteFile(filepath.Join(dir, "devcontainer-feature.json"), []byte(config), 0644)

    result, err := ParseFeatureDefinition(dir)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.ID != "node" {
        t.Errorf("expected id 'node', got %v", result.ID)
    }

    if len(result.DependsOn) != 1 || result.DependsOn[0] != "git" {
        t.Errorf("expected dependsOn [git], got %v", result.DependsOn)
    }
}
```

**Step 3: 运行测试**

Run: `go test ./pkg/config/... -v -run TestParseFeatureDefinition`
Expected: PASS

**Step 4: Commit**

```bash
git add .
git commit -m "feat: add Feature definition parsing"
```

---

### Task 8: 实现 Feature 拓扑排序

**Files:**
- Modify: `pkg/config/feature.go`

**Step 1: 添加拓扑排序函数**

```go
// pkg/config/feature.go - 添加 TopologicalSort
import "sort"

func TopologicalSort(features map[string]*FeatureDefinition) ([]string, error) {
    // Build dependency graph
    deps := make(map[string]map[string]bool)
    allDeps := make(map[string][]string)

    for id, f := range features {
        deps[id] = make(map[string]bool)
        for _, d := range f.DependsOn {
            deps[id][d] = true
        }
        allDeps[id] = f.DependsOn
    }

    // Kahn's algorithm
    inDegree := make(map[string]int)
    for id := range features {
        inDegree[id] = 0
    }
    for _, ds := range allDeps {
        for _, d := range ds {
            inDegree[d]++
        }
    }

    var queue []string
    for id, d := range inDegree {
        if d == 0 {
            queue = append(queue, id)
        }
    }

    var result []string
    for len(queue) > 0 {
        sort.Strings(queue)
        current := queue[0]
        queue = queue[1:]
        result = append(result, current)

        for id, depSet := range deps {
            if depSet[current] {
                inDegree[id]--
                if inDegree[id] == 0 {
                    queue = append(queue, id)
                }
            }
        }
    }

    if len(result) != len(features) {
        return nil, fmt.Errorf("circular dependency detected")
    }

    return result, nil
}
```

**Step 2: 添加 fmt 导入**

```go
import "fmt"
```

**Step 3: 运行测试**

Run: `go build ./pkg/config/...`
Expected: PASS

**Step 4: Commit**

```bash
git add .
git commit -m "feat: add Feature topological sorting"
```

---

### Task 9: 实现 Feature 本地读取和 OCI 拉取

**Files:**
- Create: `pkg/feature/resolver.go`
- Create: `pkg/feature/resolver_test.go`

**Step 1: 创建 Feature Resolver**

```go
// pkg/feature/resolver.go
package feature

import (
    "context"
    "fmt"
    "path/filepath"

    "github.com/devcon/cli/pkg/config"
    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

type Resolver struct {
    registryOpts []remote.Option
}

func NewResolver(registryOpts ...remote.Option) *Resolver {
    return &Resolver{
        registryOpts: registryOpts,
    }
}

func (r *Resolver) ResolveLocalFeatures(dir string) (map[string]*config.FeatureDefinition, error) {
    featuresDir := filepath.Join(dir, ".devcontainer", "features")
    // Check if directory exists
    // Parse each subdirectory's devcontainer-feature.json
    // Return map
}

func (r *Resolver) ResolveOCIFeature(ctx context.Context, ref string) (*config.FeatureDefinition, error) {
    // Parse reference like ghcr.io/devcontainers/features/node:1
    // Pull and parse devcontainer-feature.json from OCI
}
```

**Step 2: 添加依赖**

Run: `go get github.com/google/go-containerregistry`
Expected: dependency added

**Step 3: 实现本地读取**

```go
// pkg/feature/resolver.go - 实现 ResolveLocalFeatures
func (r *Resolver) ResolveLocalFeatures(dir string) (map[string]*config.FeatureDefinition, error) {
    featuresDir := filepath.Join(dir, ".devcontainer", "features")

    entries, err := os.ReadDir(featuresDir)
    if err != nil {
        if os.IsNotExist(err) {
            return make(map[string]*config.FeatureDefinition), nil
        }
        return nil, err
    }

    result := make(map[string]*config.FeatureDefinition)
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        featurePath := filepath.Join(featuresDir, entry.Name())
        def, err := config.ParseFeatureDefinition(featurePath)
        if err != nil {
            return nil, fmt.Errorf("failed to parse feature %s: %w", entry.Name(), err)
        }

        result[def.ID] = def
    }

    return result, nil
}
```

**Step 4: 实现 OCI 拉取**

```go
// pkg/feature/resolver.go - 实现 ResolveOCIFeature
func (r *Resolver) ResolveOCIFeature(ctx context.Context, ref string) (*config.FeatureDefinition, error) {
    refParsed, err := name.ParseReference(ref)
    if err != nil {
        return nil, err
    }

    img, err := remote.Image(refParsed, r.registryOpts...)
    if err != nil {
        return nil, err
    }

    // Get config
    cfg, err := img.Config()
    if err != nil {
        return nil, err
    }

    // Parse devcontainer-feature.json from label or annotation
    // ...
}
```

**Step 5: 运行测试**

Run: `go build ./pkg/feature/...`
Expected: PASS (with TODO comments)

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add Feature resolver for local and OCI"
```

---

## 阶段四：OCI 交互层

### Task 10: 实现 Feature 打包

**Files:**
- Create: `pkg/feature/package.go`
- Create: `pkg/feature/package_test.go`

**Step 1: 创建打包函数**

```go
// pkg/feature/package.go
package feature

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "os"
    "path/filepath"
)

func PackageFeature(dir string, output string) error {
    // Create tar.gz with:
    // - devcontainer-feature.json
    // - install.sh
    // - (other files)
}
```

**Step 2: 实现打包逻辑**

```go
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
```

**Step 3: 创建测试**

```go
// pkg/feature/package_test.go
package feature

import (
    "os"
    "path/filepath"
    "testing"
)

func TestPackageFeature(t *testing.T) {
    dir := t.TempDir()

    // Create test files
    os.WriteFile(filepath.Join(dir, "devcontainer-feature.json"), []byte(`{"id": "test"}`), 0644)
    os.WriteFile(filepath.Join(dir, "install.sh"), []byte("#!/bin/bash\necho hello"), 0755)

    output := filepath.Join(t.TempDir(), "output.tar.gz")

    err := PackageFeature(dir, output)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if _, err := os.Stat(output); os.IsNotExist(err) {
        t.Error("output file should exist")
    }
}
```

**Step 4: 运行测试**

Run: `go test ./pkg/feature/... -v -run TestPackageFeature`
Expected: PASS

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add Feature packaging"
```

---

### Task 11: 实现 Feature 发布

**Files:**
- Modify: `pkg/feature/package.go`

**Step 1: 添加发布函数**

```go
// pkg/feature/package.go - 添加 Publish 函数
import (
    "context"
    "github.com/google/go-containerregistry/pkg/v1/remote"
    "github.com/google/go-containerregistry/pkg/name"
)

func PublishFeature(ctx context.Context, dir string, ref string, opts ...remote.Option) error {
    // 1. Package the feature
    // 2. Parse as OCI image
    // 3. Push to registry
}
```

**Step 2: 实现发布逻辑**

```go
func PublishFeature(ctx context.Context, dir string, ref string, opts ...remote.Option) error {
    // Parse reference
    refParsed, err := name.ParseReference(ref)
    if err != nil {
        return err
    }

    // Create image from tar.gz
    // Use go-containerregistry's tarball.Image

    // Push
    if err := remote.Write(refParsed, img, opts...); err != nil {
        return err
    }

    return nil
}
```

**Step 3: Commit**

```bash
git add .
git commit -m "feat: add Feature publishing to OCI registry"
```

---

## 阶段五：构建抽象层

### Task 12: 实现 Docker Builder

**Files:**
- Create: `internal/builder/docker.go`
- Create: `internal/builder/docker_test.go`

**Step 1: 创建 Docker Builder**

```go
// internal/builder/docker.go
package builder

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/docker/docker/api/types"
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
    // 1. Prepare build context (with injected features)
    // 2. Call docker build API
    // 3. Return image ID
}
```

**Step 2: 添加 Docker 依赖**

Run: `go get github.com/docker/docker`
Expected: dependency added

**Step 3: 实现 Build 方法**

```go
// internal/builder/docker.go - 实现 Build
func (b *DockerBuilder) Build(ctx context.Context, spec Spec) (string, error) {
    // Determine build context
    contextPath := "."
    dockerfilePath := "Dockerfile"

    if spec.Dockerfile != "" {
        contextPath = filepath.Dir(spec.Dockerfile)
        dockerfilePath = filepath.Base(spec.Dockerfile)
    }

    // Read dockerfile content
    dockerfileContent, err := os.ReadFile(filepath.Join(contextPath, dockerfilePath))
    if err != nil {
        return "", err
    }

    // Prepare features injection
    // Append Feature install scripts to Dockerfile

    // Build
    resp, err := b.client.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
        Dockerfile: dockerfilePath,
        Tags:       []string{"devcon-build:latest"},
        Remove:     true,
    })
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read output
    io.Copy(os.Stdout, resp.Body)

    return "devcon-build:latest", nil
}
```

**Step 4: 实现 Up 方法**

```go
// internal/builder/docker.go - 实现 Up
func (b *DockerBuilder) Up(ctx context.Context, spec Spec) error {
    // Create and start container
    resp, err := b.client.ContainerCreate(ctx, &container.Config{
        Image: spec.Image,
        Env:   envToSlice(spec.Env),
    }, nil, nil, "")
    if err != nil {
        return err
    }

    if err := b.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
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
```

**Step 5: 编译测试**

Run: `go build ./internal/builder/...`
Expected: PASS

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add Docker builder implementation"
```

---

### Task 13: 实现 Builder Factory

**Files:**
- Create: `internal/builder/factory.go`

**Step 1: 创建 Builder Factory**

```go
// internal/builder/factory.go
package builder

import (
    "context"
    "fmt"
)

func NewBuilder(provider string) (Builder, error) {
    switch provider {
    case "docker":
        return NewDockerBuilder()
    // case "buildkit":
    //     return NewBuildKitBuilder()
    // case "kaniko":
    //     return NewKanikoBuilder()
    default:
        return nil, fmt.Errorf("unknown provider: %s", provider)
    }
}
```

**Step 2: Commit**

```bash
git add .
git commit -m "feat: add Builder factory"
```

---

## 阶段六：CLI 命令实现

### Task 14: 实现 features package 命令

**Files:**
- Modify: `cmd/devcon/features_package.go`

**Step 1: 实现 package 命令**

```go
// cmd/devcon/features_package.go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/devcon/cli/pkg/feature"
    "github.com/spf13/cobra"
)

var featuresPackageCmd = &cobra.Command{
    Use:   "package <dir>",
    Short: "Package a Feature as OCI artifact",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]
        output, _ := cmd.Flags().GetString("output")

        if output == "" {
            output = filepath.Join(dir, "feature.tar.gz")
        }

        if err := feature.PackageFeature(dir, output); err != nil {
            return err
        }

        fmt.Println("Feature packaged:", output)
        return nil
    },
}

func init() {
    featuresPackageCmd.Flags().String("output", "", "Output path")
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: 测试命令**

```bash
# Create test feature
mkdir -p /tmp/test-feature
echo '{"id": "test", "name": "Test"}' > /tmp/test-feature/devcontainer-feature.json
echo '#!/bin/bash' > /tmp/test-feature/install.sh

# Run package
./devcon features package /tmp/test-feature
```

**Step 4: Commit**

```bash
git add .
git commit -m "feat: implement features package command"
```

---

### Task 15: 实现 features publish 命令

**Files:**
- Modify: `cmd/devcon/features_publish.go`

**Step 1: 实现 publish 命令**

```go
// cmd/devcon/features_publish.go
package cmd

import (
    "context"
    "fmt"

    "github.com/devcon/cli/pkg/feature"
    "github.com/spf13/cobra"
)

var featuresPublishCmd = &cobra.Command{
    Use:   "publish <dir>",
    Short: "Publish a Feature to OCI registry",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]
        reg, _ := cmd.Flags().GetString("reg")

        if reg == "" {
            return fmt.Errorf("--reg is required")
        }

        ctx := context.Background()
        if err := feature.PublishFeature(ctx, dir, reg); err != nil {
            return err
        }

        fmt.Println("Feature published:", reg)
        return nil
    },
}

func init() {
    featuresPublishCmd.Flags().String("reg", "", "Registry URL (e.g., ghcr.io/user/feature)")
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement features publish command"
```

---

### Task 16: 实现 build 命令

**Files:**
- Modify: `cmd/devcon/build.go`

**Step 1: 实现 build 命令**

```go
// cmd/devcon/build.go
package cmd

import (
    "context"
    "fmt"

    "github.com/devcon/cli/internal/builder"
    "github.com/devcon/cli/pkg/config"
    "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
    Use:   "build <dir>",
    Short: "Build devcontainer image",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]
        provider, _ := cmd.Flags().GetString("provider")

        // Parse config
        cfg, err := config.ParseDevcontainer(dir)
        if err != nil {
            return fmt.Errorf("failed to parse config: %w", err)
        }

        // Resolve extends
        cfg, err = config.ResolveExtends(dir, cfg)
        if err != nil {
            return fmt.Errorf("failed to resolve extends: %w", err)
        }

        // Create builder
        b, err := builder.NewBuilder(provider)
        if err != nil {
            return err
        }

        // Build
        spec := builder.Spec{
            Image:      cfg.Image,
            Dockerfile: cfg.Dockerfile,
            Features:   cfg.Features,
            Env:        cfg.Env,
        }

        imageID, err := b.Build(context.Background(), spec)
        if err != nil {
            return err
        }

        fmt.Println("Image built:", imageID)
        return nil
    },
}

func init() {
    buildCmd.Flags().String("provider", "docker", "Build provider (docker)")
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement build command"
```

---

### Task 17: 实现 config 命令

**Files:**
- Modify: `cmd/devcon/config.go`

**Step 1: 实现 config 命令**

```go
// cmd/devcon/config.go
package cmd

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/devcon/cli/pkg/config"
    "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
    Use:   "config <dir>",
    Short: "Validate and show devcontainer config",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]

        cfg, err := config.ParseDevcontainer(dir)
        if err != nil {
            return fmt.Errorf("failed to parse config: %w", err)
        }

        cfg, err = config.ResolveExtends(dir, cfg)
        if err != nil {
            return fmt.Errorf("failed to resolve extends: %w", err)
        }

        output, _ := cmd.Flags().GetString("output")
        if output == "json" {
            encoder := json.NewEncoder(os.Stdout)
            encoder.SetIndent("", "  ")
            encoder.Encode(cfg)
        } else {
            fmt.Printf("Image: %s\n", cfg.Image)
            fmt.Printf("Dockerfile: %s\n", cfg.Dockerfile)
            fmt.Printf("Features: %v\n", cfg.Features)
        }

        return nil
    },
}

func init() {
    configCmd.Flags().String("output", "text", "Output format (text, json)")
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement config command"
```

---

### Task 18: 实现 inspect 命令

**Files:**
- Modify: `cmd/devcon/inspect.go`

**Step 1: 实现 inspect 命令**

```go
// cmd/devcon/inspect.go
package cmd

import (
    "context"
    "fmt"

    "github.com/devcon/cli/pkg/config"
    "github.com/devcon/cli/pkg/feature"
    "github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
    Use:   "inspect <dir>",
    Short: "Inspect parsed config and feature dependencies",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]

        // Parse config
        cfg, err := config.ParseDevcontainer(dir)
        if err != nil {
            return fmt.Errorf("failed to parse config: %w", err)
        }

        cfg, err = config.ResolveExtends(dir, cfg)
        if err != nil {
            return fmt.Errorf("failed to resolve extends: %w", err)
        }

        // Resolve features
        resolver := feature.NewResolver()

        localFeatures, err := resolver.ResolveLocalFeatures(dir)
        if err != nil {
            return fmt.Errorf("failed to resolve local features: %w", err)
        }

        fmt.Println("=== Config ===")
        fmt.Printf("Image: %s\n", cfg.Image)
        fmt.Printf("Features: %v\n", cfg.Features)

        fmt.Println("\n=== Local Features ===")
        for id, f := range localFeatures {
            fmt.Printf("- %s (version: %s)\n", id, f.Version)
            if len(f.DependsOn) > 0 {
                fmt.Printf("  dependsOn: %v\n", f.DependsOn)
            }
        }

        return nil
    },
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement inspect command"
```

---

### Task 19: 实现 up 命令

**Files:**
- Modify: `cmd/devcon/up.go`

**Step 1: 实现 up 命令**

```go
// cmd/devcon/up.go
package cmd

import (
    "context"
    "fmt"

    "github.com/devcon/cli/internal/builder"
    "github.com/devcon/cli/pkg/config"
    "github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
    Use:   "up <dir>",
    Short: "Start devcontainer for local testing",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := args[0]
        provider, _ := cmd.Flags().GetString("provider")

        // Parse config
        cfg, err := config.ParseDevcontainer(dir)
        if err != nil {
            return fmt.Errorf("failed to parse config: %w", err)
        }

        cfg, err = config.ResolveExtends(dir, cfg)
        if err != nil {
            return fmt.Errorf("failed to resolve extends: %w", err)
        }

        // Build first if needed
        b, err := builder.NewBuilder(provider)
        if err != nil {
            return err
        }

        spec := builder.Spec{
            Image:      cfg.Image,
            Dockerfile: cfg.Dockerfile,
            Features:   cfg.Features,
            Env:        cfg.Env,
        }

        imageID, err := b.Build(context.Background(), spec)
        if err != nil {
            return err
        }

        // Start container
        spec.Image = imageID
        if err := b.Up(context.Background(), spec); err != nil {
            return err
        }

        return nil
    },
}

func init() {
    upCmd.Flags().String("provider", "docker", "Build provider (docker)")
}
```

**Step 2: 编译测试**

Run: `go build -o devcon ./cmd/devcon`
Expected: PASS

**Step 3: Commit**

```bash
git add .
git commit -m "feat: implement up command"
```

---

## 总结

实现顺序：
1. 项目初始化
2. 配置解析层
3. Feature 调度层
4. OCI 交互层
5. 构建抽象层
6. CLI 命令实现
