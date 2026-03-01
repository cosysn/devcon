package builder

import (
	"context"
)

type Spec struct {
	ContextDir         string
	Image              string
	Dockerfile         string
	Features          map[string]interface{}
	Env                map[string]string
	Mounts             []string
	Ports              []int
	RemoteUser         string
	Privileged         bool
	OnCreateCommand    string
	PostCreateCommand  string
	PostStartCommand   string
}

type Builder interface {
	Build(ctx context.Context, spec Spec) (string, error)
	Up(ctx context.Context, spec Spec) error
}
