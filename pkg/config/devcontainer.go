package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StringOrArray represents a value that can be either a string or an array of strings
type StringOrArray []string

// UnmarshalJSON implements custom unmarshaling for StringOrArray
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	// Try as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = []string{str}
		return nil
	}

	// Try as array
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = arr
		return nil
	}

	return fmt.Errorf("invalid StringOrArray: %s", string(data))
}

// MarshalJSON implements custom marshaling for StringOrArray
func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}
	return json.Marshal([]string(s))
}

// ToSlice converts StringOrArray to a string slice
func (s StringOrArray) ToSlice() []string {
	return []string(s)
}

// ToString converts StringOrArray to a single string (joins with &&)
func (s StringOrArray) ToString() string {
	return strings.Join(s, " && ")
}

// DevcontainerConfig represents the devcontainer.json configuration
type DevcontainerConfig struct {
	Name                string                 `json:"name,omitempty"`
	Image               string                 `json:"image,omitempty"`
	Dockerfile          string                 `json:"dockerFile,omitempty"`
	Build              *BuildConfig           `json:"build,omitempty"`
	Features                  map[string]interface{} `json:"features,omitempty"`
	OverrideFeatureInstallOrder []string               `json:"overrideFeatureInstallOrder,omitempty"`
	ContainerEnv       map[string]string      `json:"containerEnv,omitempty"`
	RemoteEnv          map[string]string      `json:"remoteEnv,omitempty"`
	Mounts             []string               `json:"mounts,omitempty"`
	Ports              []int                  `json:"forwardPorts,omitempty"`
	PortsAttributes    map[string]PortAttribute `json:"portsAttributes,omitempty"`
	OtherPortsAttributes *PortAttribute        `json:"otherPortsAttributes,omitempty"`
	RemoteUser         string                 `json:"remoteUser,omitempty"`
	ContainerUser      string                 `json:"containerUser,omitempty"`
	WorkspaceMount     string                 `json:"workspaceMount,omitempty"`
	WorkspaceFolder    string                 `json:"workspaceFolder,omitempty"`
	Extends            string                 `json:"extends,omitempty"`
	InitializeCommand  StringOrArray       `json:"initializeCommand,omitempty"`
	OnCreateCommand    StringOrArray       `json:"onCreateCommand,omitempty"`
	UpdateContentCommand StringOrArray      `json:"updateContentCommand,omitempty"`
	PostCreateCommand  StringOrArray       `json:"postCreateCommand,omitempty"`
	PostStartCommand   StringOrArray       `json:"postStartCommand,omitempty"`
	PostAttachCommand  StringOrArray       `json:"postAttachCommand,omitempty"`
	WaitFor           string                 `json:"waitFor,omitempty"`
	OverrideCommand    *bool                  `json:"overrideCommand,omitempty"`
	ShutdownAction     string                 `json:"shutdownAction,omitempty"`
	HostRequirements  *HostRequirements      `json:"hostRequirements,omitempty"`
	RunArgs           []string               `json:"runArgs,omitempty"`
	UpdateRemoteUserUID *bool                `json:"updateRemoteUserUID,omitempty"`
	UserEnvProbe      string                 `json:"userEnvProbe,omitempty"`
	Init              *bool                  `json:"init,omitempty"`
	Privileged        *bool                  `json:"privileged,omitempty"`
	CapAdd            []string               `json:"capAdd,omitempty"`
	SecurityOpt       []string               `json:"securityOpt,omitempty"`
	DockerComposeFile interface{}            `json:"dockerComposeFile,omitempty"`
	Service           string                 `json:"service,omitempty"`
	RunServices       []string               `json:"runServices,omitempty"`
	AppPort           interface{}            `json:"appPort,omitempty"`
	Customizations    *Customizations        `json:"customizations,omitempty"`
}

type BuildConfig struct {
	Dockerfile string            `json:"dockerfile,omitempty"`
	Context    string            `json:"context,omitempty"`
	Args       map[string]string `json:"args,omitempty"`
	Options    []string          `json:"options,omitempty"`
	Target     string            `json:"target,omitempty"`
	CacheFrom  interface{}       `json:"cacheFrom,omitempty"`
}

// HostRequirements specifies minimum host hardware requirements
type HostRequirements struct {
	CPUs    int    `json:"cpus,omitempty"`
	Memory  string `json:"memory,omitempty"`
	Storage string `json:"storage,omitempty"`
}

// PortAttribute defines attributes for a forwarded port
type PortAttribute struct {
	Label           string `json:"label,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	OnAutoForward   string `json:"onAutoForward,omitempty"`
	RequireLocalPort *bool `json:"requireLocalPort,omitempty"`
	ElevateIfNeeded *bool  `json:"elevateIfNeeded,omitempty"`
}

type Customizations struct {
	VSCode *VSCodeCustomization `json:"vscode,omitempty"`
}

type VSCodeCustomization struct {
	Extensions []string               `json:"extensions,omitempty"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
}

func ParseDevcontainer(dir string) (*DevcontainerConfig, error) {
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read devcontainer.json: %w", err)
	}

	parsed, err := ParseJSONC(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSONC: %w", err)
	}

	// Simple JSON unmarshal into struct
	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parsed config: %w", err)
	}
	var config DevcontainerConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// If Dockerfile is specified in build config, promote it to top-level
	if config.Dockerfile == "" && config.Build != nil && config.Build.Dockerfile != "" {
		config.Dockerfile = config.Build.Dockerfile
	}

	// Resolve devcontainer variables (e.g., ${localWorkspaceFolder})
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}
	resolveVariables(&config, absDir)

	return &config, nil
}

// resolveVariables replaces devcontainer variables with their values
func resolveVariables(config *DevcontainerConfig, localWorkspaceFolder string) {
	// Resolve mounts
	if config.Mounts != nil {
		resolved := make([]string, len(config.Mounts))
		for i, m := range config.Mounts {
			resolved[i] = strings.ReplaceAll(m, "${localWorkspaceFolder}", localWorkspaceFolder)
		}
		config.Mounts = resolved
	}

	// Resolve workspaceMount
	if config.WorkspaceMount != "" {
		config.WorkspaceMount = strings.ReplaceAll(config.WorkspaceMount, "${localWorkspaceFolder}", localWorkspaceFolder)
	}
}

// ResolveExtends resolves the extends property in devcontainer config
// It recursively merges the base config with the current config
func ResolveExtends(dir string, config *DevcontainerConfig) (*DevcontainerConfig, error) {
	if config.Extends == "" {
		return config, nil
	}

	// Resolve extends path - could be "./base.json" or just "base"
	extendsPath := config.Extends
	if !filepath.IsAbs(extendsPath) && !strings.HasPrefix(extendsPath, "./") && !strings.HasPrefix(extendsPath, "../") {
		extendsPath = "./" + extendsPath
	}

	fullPath := filepath.Join(dir, ".devcontainer", extendsPath)
	// Add .json if not present
	if !strings.HasSuffix(fullPath, ".json") {
		fullPath += ".json"
	}

	// Clean and verify path stays within .devcontainer
	cleanPath := filepath.Clean(fullPath)
	devcontainerDir := filepath.Clean(filepath.Join(dir, ".devcontainer"))
	if !strings.HasPrefix(cleanPath, devcontainerDir) {
		return nil, fmt.Errorf("extends path escapes devcontainer directory: %s", fullPath)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read extended config %s: %w", fullPath, err)
	}

	parsed, err := ParseJSONC(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse extended config: %w", err)
	}

	var base DevcontainerConfig
	jsonData, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal base config: %w", err)
	}
	if err := json.Unmarshal(jsonData, &base); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base config: %w", err)
	}

	// Recursively resolve extends
	resolvedBase, err := ResolveExtends(dir, &base)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve base extends: %w", err)
	}
	base = *resolvedBase

	// If Dockerfile is specified in build config, promote it to top-level
	if base.Dockerfile == "" && base.Build != nil && base.Build.Dockerfile != "" {
		base.Dockerfile = base.Build.Dockerfile
	}

	// Merge: base -> config (config takes precedence)
	if config.Name != "" {
		base.Name = config.Name
	}
	if config.Image != "" {
		base.Image = config.Image
	}
	if config.Dockerfile != "" {
		base.Dockerfile = config.Dockerfile
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

	// Merge containerEnv
	if config.ContainerEnv != nil {
		if base.ContainerEnv == nil {
			base.ContainerEnv = make(map[string]string)
		}
		for k, v := range config.ContainerEnv {
			base.ContainerEnv[k] = v
		}
	}

	// Merge remoteEnv
	if config.RemoteEnv != nil {
		if base.RemoteEnv == nil {
			base.RemoteEnv = make(map[string]string)
		}
		for k, v := range config.RemoteEnv {
			base.RemoteEnv[k] = v
		}
	}

	// Merge mounts
	if config.Mounts != nil {
		base.Mounts = append(base.Mounts, config.Mounts...)
	}

	// Merge ports
	if config.Ports != nil {
		base.Ports = append(base.Ports, config.Ports...)
	}

	// Merge portsAttributes
	if config.PortsAttributes != nil {
		if base.PortsAttributes == nil {
			base.PortsAttributes = make(map[string]PortAttribute)
		}
		for k, v := range config.PortsAttributes {
			base.PortsAttributes[k] = v
		}
	}

	// Merge otherPortsAttributes
	if config.OtherPortsAttributes != nil {
		base.OtherPortsAttributes = config.OtherPortsAttributes
	}

	// Merge user-related fields
	if config.RemoteUser != "" {
		base.RemoteUser = config.RemoteUser
	}
	if config.ContainerUser != "" {
		base.ContainerUser = config.ContainerUser
	}
	if config.WorkspaceMount != "" {
		base.WorkspaceMount = config.WorkspaceMount
	}
	if config.WorkspaceFolder != "" {
		base.WorkspaceFolder = config.WorkspaceFolder
	}
	if config.Extends != "" {
		base.Extends = config.Extends
	}

	// Merge lifecycle commands (array concatenation)
	if len(config.InitializeCommand) > 0 {
		base.InitializeCommand = append(base.InitializeCommand, config.InitializeCommand.ToSlice()...)
	}
	if len(config.OnCreateCommand) > 0 {
		base.OnCreateCommand = append(base.OnCreateCommand, config.OnCreateCommand.ToSlice()...)
	}
	if len(config.UpdateContentCommand) > 0 {
		base.UpdateContentCommand = append(base.UpdateContentCommand, config.UpdateContentCommand.ToSlice()...)
	}
	if len(config.PostCreateCommand) > 0 {
		base.PostCreateCommand = append(base.PostCreateCommand, config.PostCreateCommand.ToSlice()...)
	}
	if len(config.PostStartCommand) > 0 {
		base.PostStartCommand = append(base.PostStartCommand, config.PostStartCommand.ToSlice()...)
	}
	if len(config.PostAttachCommand) > 0 {
		base.PostAttachCommand = append(base.PostAttachCommand, config.PostAttachCommand.ToSlice()...)
	}

	// Merge command handling fields
	if config.WaitFor != "" {
		base.WaitFor = config.WaitFor
	}
	if config.OverrideCommand != nil {
		base.OverrideCommand = config.OverrideCommand
	}
	if config.ShutdownAction != "" {
		base.ShutdownAction = config.ShutdownAction
	}
	if config.AppPort != nil {
		base.AppPort = config.AppPort
	}

	// Merge host requirements
	if config.HostRequirements != nil {
		if base.HostRequirements == nil {
			base.HostRequirements = &HostRequirements{}
		}
		if config.HostRequirements.CPUs > 0 {
			base.HostRequirements.CPUs = config.HostRequirements.CPUs
		}
		if config.HostRequirements.Memory != "" {
			base.HostRequirements.Memory = config.HostRequirements.Memory
		}
		if config.HostRequirements.Storage != "" {
			base.HostRequirements.Storage = config.HostRequirements.Storage
		}
	}

	// Merge security options (individual fields)
	if config.Init != nil {
		base.Init = config.Init
	}
	if config.Privileged != nil {
		base.Privileged = config.Privileged
	}
	if len(config.CapAdd) > 0 {
		base.CapAdd = append(base.CapAdd, config.CapAdd...)
	}
	if len(config.SecurityOpt) > 0 {
		base.SecurityOpt = append(base.SecurityOpt, config.SecurityOpt...)
	}

	// Merge runArgs
	if len(config.RunArgs) > 0 {
		base.RunArgs = append(base.RunArgs, config.RunArgs...)
	}

	// Merge user environment probe and UID update
	if config.UpdateRemoteUserUID != nil {
		base.UpdateRemoteUserUID = config.UpdateRemoteUserUID
	}
	if config.UserEnvProbe != "" {
		base.UserEnvProbe = config.UserEnvProbe
	}
	if config.Init != nil {
		base.Init = config.Init
	}
	if config.Privileged != nil {
		base.Privileged = config.Privileged
	}

	// Merge Docker Compose config
	if config.DockerComposeFile != nil {
		base.DockerComposeFile = config.DockerComposeFile
	}
	if config.Service != "" {
		base.Service = config.Service
	}
	if len(config.RunServices) > 0 {
		base.RunServices = append(base.RunServices, config.RunServices...)
	}

	// Merge customizations
	if config.Customizations != nil {
		if base.Customizations == nil {
			base.Customizations = &Customizations{}
		}
		if config.Customizations.VSCode != nil {
			if base.Customizations.VSCode == nil {
				base.Customizations.VSCode = &VSCodeCustomization{}
			}
			// Merge extensions
			if len(config.Customizations.VSCode.Extensions) > 0 {
				base.Customizations.VSCode.Extensions = append(base.Customizations.VSCode.Extensions, config.Customizations.VSCode.Extensions...)
			}
			// Merge settings
			if config.Customizations.VSCode.Settings != nil {
				if base.Customizations.VSCode.Settings == nil {
					base.Customizations.VSCode.Settings = make(map[string]interface{})
				}
				for k, v := range config.Customizations.VSCode.Settings {
					base.Customizations.VSCode.Settings[k] = v
				}
			}
		}
	}

	return &base, nil
}
