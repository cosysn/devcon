package config

// DevcontainerConfig represents the devcontainer.json configuration
type DevcontainerConfig struct {
    Image         string                 `json:"image,omitempty"`
    Build         *BuildConfig           `json:"build,omitempty"`
    Features      map[string]interface{} `json:"features,omitempty"`
    Env           map[string]string      `json:"env,omitempty"`
    Mounts        []string               `json:"mounts,omitempty"`
    Ports         []int                  `json:"ports,omitempty"`
    RemoteUser    string                 `json:"remoteUser,omitempty"`
}

type BuildConfig struct {
    Dockerfile string `json:"dockerfile,omitempty"`
    Context    string `json:"context,omitempty"`
}
